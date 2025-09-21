package entity

import (
	"encore.dev/types/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Name   string    `json:"name"`
	AppIds []string  `json:"appIds"`
	PIID   uuid.UUID `json:"piid"`
	jwt.RegisteredClaims
}
type AuthData struct {
	AppMapping map[string]bool
	PIID       uuid.UUID
}
