package internal

import (
	"time"

	"encore.app/authentication/entity"
	"encore.dev/types/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type (
	ITokenRepo interface {
		GenerateToken(username string, userId uuid.UUID, piidApps map[uuid.UUID]entity.ProductAndApps, expirationDuration time.Duration) *jwt.Token
		TokenToSignedString(jwtToken *jwt.Token) (string, error)
	}

	ITokenGenerator interface {
		GenerateToken(userName string, userId uuid.UUID, piidApps map[uuid.UUID]entity.ProductAndApps, expirationTime time.Duration) (string, error)
	}
)

type TokenGenerator struct {
	r ITokenRepo
}

func NewTokenGenerator(r ITokenRepo) TokenGenerator {
	return TokenGenerator{r: r}
}

func (g TokenGenerator) GenerateToken(userName string, userId uuid.UUID, piidApps map[uuid.UUID]entity.ProductAndApps, expirationTime time.Duration) (string, error) {
	token := g.r.GenerateToken(userName, userId, piidApps, expirationTime)
	return g.r.TokenToSignedString(token)
}
