package service

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/open-telemetry/opamp-go/opamp-server/data"
	"github.com/open-telemetry/opamp-go/opamp-server/login/models"
	"github.com/open-telemetry/opamp-go/protobufs"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ClientID = "open-observability"

	Issuer    = "http://localhost:8080/realms/global"
	JWKSURL   = Issuer + "/protocol/openid-connect/certs"
	ClaimsKey = "claims"
)

var (
	cachedJWKS *models.JWKS
	lastFetch  time.Time
	mu         sync.Mutex
)

//-----------------------------------------------------Validation----------------------------------------------

func ValidateToken(tokenString string) (*models.IDPClaims, error) {

	// Bearer parsing
	tokenString = strings.TrimSpace(tokenString)
	if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
		tokenString = tokenString[7:]
	} else {
		return nil, errors.New("invalid authorization header")
	}

	jwks, err := getJWKS()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&models.IDPClaims{},
		func(t *jwt.Token) (any, error) {

			if t.Method.Alg() != "RS256" {
				return nil, errors.New("invalid signing method")
			}

			kid, ok := t.Header["kid"].(string)
			if !ok {
				return nil, errors.New("missing kid")
			}

			return getPublicKey(jwks, kid)
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.IDPClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// ISSUER CHECK
	if claims.Issuer != Issuer {
		return nil, errors.New("invalid issuer")
	}

	// EXPIRY CHECK
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func getJWKS() (*models.JWKS, error) {

	if cachedJWKS != nil && time.Since(lastFetch) < time.Hour {
		return cachedJWKS, nil
	}

	mu.Lock()
	defer mu.Unlock()

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(JWKSURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch JWKS: status " + http.StatusText(resp.StatusCode))
	}

	var jwks models.JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	if len(jwks.Keys) == 0 {
		return nil, errors.New("empty JWKS keys")
	}

	cachedJWKS = &jwks
	lastFetch = time.Now()

	return cachedJWKS, nil
}

func getPublicKey(jwks *models.JWKS, kid string) (*rsa.PublicKey, error) {

	for _, k := range jwks.Keys {
		if k.Kid != kid {
			continue
		}

		// decode modulus
		nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
		if err != nil {
			return nil, err
		}

		// decode exponent (SAFE)
		eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
		if err != nil {
			return nil, err
		}

		e := 0
		for _, b := range eBytes {
			e = e<<8 | int(b)
		}

		return &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: e,
		}, nil
	}

	return nil, errors.New("public key not found for kid")
}

//-----------------------------------------Roles Extraction--------------------------------------------------------

func ExtractRolesAndPermissions(claims *models.IDPClaims) (roles []string, permissions []string) {

	roles = claims.RealmAccess.Roles

	if client, ok := claims.ResourceAccess[ClientID]; ok {
		permissions = client.Roles
	}

	return roles, permissions
}

//-------------------------------------------Agent Info Extraction--------------------------------------------------

func ExtractAgentInfo(agent data.Agent) (models.AgentResponse, error) {
	info := models.AgentResponse{}

	if agent.Status.AgentDescription != nil {
		attrs := agent.Status.AgentDescription.IdentifyingAttributes
		info.ServerName = getAttribute(attrs, "service.name")
		info.SupervisorVersion = getAttribute(attrs, "service.version")
	}

	if agent.Status.AgentDescription != nil {
		attrs := agent.Status.AgentDescription.NonIdentifyingAttributes
		info.AgentVersion = getAttribute(attrs, "host.name")
		info.OS = getAttribute(attrs, "os.type")
	}

	if agent.Status.Health != nil {
		info.Active = agent.Status.Health.Healthy
	}

	return info, nil
}

func getAttribute(attrs []*protobufs.KeyValue, key string) string {
	for _, kv := range attrs {
		if kv.Key == key {
			return kv.Value.GetStringValue()
		}
	}
	return ""
}
