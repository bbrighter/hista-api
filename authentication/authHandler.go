package authentication

import (
	"context"
	"crypto/rsa"
	"fmt"

	"encore.app/entity"
	"encore.app/errors"
	"encore.dev/beta/auth"
	"github.com/golang-jwt/jwt/v5"
)

type AuthParams struct {
	Token string `header:"Authorization"`
}

//encore:authhandler
func AuthHandler(ctx context.Context, params *AuthParams) (auth.UID, *entity.AuthData, error) {
	token, err := parseToken(params.Token)
	if err != nil {
		return auth.UID(""), &entity.AuthData{}, errors.ErrorUnauthenticated
	}
	var appMapping = make(map[string]bool)

	authData := &entity.AuthData{
		AppMapping: appMapping,
		PIID:       token.PIID,
	}

	return auth.UID(token.Subject), authData, nil
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
