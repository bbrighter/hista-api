package authentication

import (
	"context"
	"crypto/rsa"
	"fmt"

	"encore.app/authentication/entity"
	"encore.app/errors"
	shared "encore.app/shared/entity"
	"encore.dev/beta/auth"
	"github.com/golang-jwt/jwt/v5"
)

type AuthParams struct {
	Token string `header:"Authorization"`
}

//encore:authhandler
func AuthHandler(ctx context.Context, params *AuthParams) (auth.UID, *shared.AuthData, error) {
	claims, err := parseToken(params.Token)
	if err != nil {
		return auth.UID(""), &shared.AuthData{}, errors.ErrorUnauthenticated
	}

	var instances []shared.AuthProductInstance
	for _, c := range claims.ProductInstances {
		var appMapping = make(map[string]bool)
		for _, a := range c.AppIds {
			appMapping[a] = true
		}
		inst := shared.AuthProductInstance{PIID: c.PIID, AppMapping: appMapping}
		instances = append(instances, inst)
	}
	return auth.UID(claims.Subject), &shared.AuthData{Instances: instances}, nil
}

func parseToken(token string) (*entity.CustomClaims, error) {
	jwtToken, err := validateToken(token, keys.publicKey)
	if err != nil {
		return nil, err
	}
	customClaims, ok := jwtToken.Claims.(*entity.CustomClaims)
	if !ok {
		return nil, errors.ErrorUnauthenticated
	}
	return customClaims, nil
}

func validateToken(token string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	var claims entity.CustomClaims
	jwtToken, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return publicKey, nil
		},
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return jwtToken, err
	}
	if !jwtToken.Valid {
		return jwtToken, errors.ErrorUnauthenticated
	}
	return jwtToken, nil
}
