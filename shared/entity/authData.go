package entity

import (
	"encore.dev/types/uuid"
)

type AuthData struct {
	Instances []AuthProductInstance `json:"instances"`
	UserName  string                `json:"userName"`
	UserId    uuid.UUID             `json:"userId"`
}

type AuthProductInstance struct {
	AppMapping map[string]bool `json:"appMapping"`
	PIID       uuid.UUID       `json:"piid"`
	Product    string          `json:"product"`
}
