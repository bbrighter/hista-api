package internal

import (
	"time"

	"encore.dev/types/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type (
	ITokenRepo interface {
		GenerateToken(username string, userId uuid.UUID, appIds []string, piid uuid.UUID, expirationDuration time.Duration) *jwt.Token
		TokenToSignedString(jwtToken *jwt.Token) (string, error)
	}

	ITokenGenerator interface {
		GenerateToken(userName string, userId uuid.UUID, appIds []string, piid uuid.UUID) (string, error)
	}
)

type TokenGenerator struct {
	r ITokenRepo
}

func NewTokenGenerator(r ITokenRepo) TokenGenerator {
	return TokenGenerator{r: r}
}

func (g TokenGenerator) GenerateToken(userName string, userId uuid.UUID, appIds []string, piid uuid.UUID) (string, error) {
	expirationTime := time.Hour * 24
	token := g.r.GenerateToken(userName, userId, appIds, piid, expirationTime)
	return g.r.TokenToSignedString(token)
}
