package authentication

import (
	"crypto/rsa"
	"fmt"
	"time"

	"encore.app/errors"
	"encore.dev/types/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Name   string    `json:"name"`
	AppIds []string  `json:"appIds"`
	PIID   uuid.UUID `json:"piid"`
	jwt.RegisteredClaims
}

func GenerateToken(username string, userId uuid.UUID, appIds []string, piid uuid.UUID, expirationDuration time.Duration) *jwt.Token {
	claims := CustomClaims{
		Name:   username,
		AppIds: appIds,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationDuration)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
}

// func tokenToSignedString(jwtToken *jwt.Token, privateKey *rsa.PrivateKey) (string, error) {
// 	return jwtToken.SignedString(privateKey)
// }

func validateToken(token string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	var claims CustomClaims
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

func parseToken(token string, publicKey *rsa.PublicKey) (*CustomClaims, error) {
	jwtToken, err := validateToken(token, publicKey)
	if err != nil {
		return nil, err
	}
	customClaims, ok := jwtToken.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.ErrorUnauthenticated
	}
	return customClaims, nil
}
