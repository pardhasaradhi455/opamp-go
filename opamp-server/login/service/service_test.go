package service

import "testing"

func TestGetJwks(t *testing.T) {
	jwks, err := getJWKS()

	if err != nil {
		t.Error(err.Error())
	}

	if jwks == nil || len(jwks.Keys) == 0 {
		t.Error("jwks empty")
	}
}

func TestGetPublicKey(t *testing.T) {
	jwks, err := getJWKS()

	if err != nil {
		t.Error(err.Error())
	}

	kid := jwks.Keys[0].Kid

	publicKey, err := getPublicKey(jwks, kid)

	if err != nil {
		t.Error(err.Error())
	}

	if publicKey == nil {
		t.Error("Empty public key")
	}
}
