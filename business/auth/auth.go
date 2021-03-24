// Package auth provides authentication and authorization support.
package auth

import (
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go/v4"
)

// These are the expected values for Claims.Roles.
const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// Key is used to store/retrieve a Claims value from a context.Context.
const Key ctxKey = 1

// Claims represents the authorization claims transmitted via a JWT.
type Claims struct {
	jwt.StandardClaims
	Roles []string `json:"roles"`
}

// HasRole returns true if the claims has at least one of the provided roles.
func (c Claims) HasRole(roles ...string) bool {
	for _, has := range c.Roles {
		for _, want := range roles {
			if has == want {
				return true
			}
		}
	}
	return false
}

// Keys represents an in memory store of keys.
type Keys map[string]*rsa.PrivateKey

/*
PublicKeyLookup defines the signature of a function to lookup public keys.

In a production system, a key id (KID) is used to retrieve the correct public key to parse a JWT for auth and claims.
A key lookup function is provided to perform the task of retrieving a KID for a given public key.

A key lookup function is required for creating an Authenticator.

 * Private keys should be rotated. During the transition period, tokens signed with the old and new keys can coexist by looking up the correct public key by KID.

 * KID to public key resolution is usually accomplished via a public JWKS endpoint.
 See https://auth0.com/docs/jwks for more details.
*/
type PublicKeyLookup func(kid string) (*rsa.PublicKey, err)
