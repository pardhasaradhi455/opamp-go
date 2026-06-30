package models

import "github.com/golang-jwt/jwt/v5"

type IDPClaims struct {
	PreferredUsername string `json:"preferred_username"`
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`

	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`

	jwt.RegisteredClaims
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	// General fields
	Kid string `json:"kid,omitempty"`
	Kty string `json:"kty,omitempty"`
	Use string `json:"use,omitempty"`
	Alg string `json:"alg,omitempty"`

	// RSA fields
	N  string `json:"n,omitempty"`  // Modulus
	E  string `json:"e,omitempty"`  // Exponent
	D  string `json:"d,omitempty"`  // Private exponent
	P  string `json:"p,omitempty"`  // First prime factor
	Q  string `json:"q,omitempty"`  // Second prime factor
	DP string `json:"dp,omitempty"` // First factor’s CRT exponent
	DQ string `json:"dq,omitempty"` // Second factor’s CRT exponent
	QI string `json:"qi,omitempty"` // CRT coefficient

	// EC fields
	Crv string `json:"crv,omitempty"` // Curve name
	X   string `json:"x,omitempty"`   // X coordinate
	Y   string `json:"y,omitempty"`   // Y coordinate

	// Symmetric key field
	K string `json:"k,omitempty"` // Symmetric key value

	// X.509 certificate fields
	X5c     []string `json:"x5c,omitempty"`      // Certificate chain
	X5t     string   `json:"x5t,omitempty"`      // SHA-1 thumbprint
	X5tS256 string   `json:"x5t#S256,omitempty"` // SHA-256 thumbprint
}
