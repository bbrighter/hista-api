package repository

import (
	"crypto/rsa"
	"time"

	"encore.app/authentication/entity"
	"encore.dev/types/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type TokenRepo struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewTokenRepo(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *TokenRepo {
	return &TokenRepo{privateKey: privateKey, publicKey: publicKey}
}

func (r *TokenRepo) GenerateToken(username string, userId uuid.UUID, piidApps map[uuid.UUID][]string, expirationDuration time.Duration) *jwt.Token {
	var instances []entity.ProductInstanceClaim
	for piid, appIds := range piidApps {
		instances = append(instances, entity.ProductInstanceClaim{PIID: piid, AppIds: appIds})
	}
	claims := entity.CustomClaims{
		UserName:         username,
		ProductInstances: instances,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationDuration)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
}

func (r *TokenRepo) TokenToSignedString(jwtToken *jwt.Token) (string, error) {
	return jwtToken.SignedString(r.privateKey)
}
