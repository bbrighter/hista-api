package entity

import (
	"encore.dev/types/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserName         string                 `json:"userName"`
	ProductInstances []ProductInstanceClaim `json:"productInstances"`
	jwt.RegisteredClaims
}

type ProductInstanceClaim struct {
	PIID   uuid.UUID `json:"piid"`
	AppIds []string  `json:"appIds"`
}
