package entity

import (
	"encore.dev/types/uuid"
)

type AuthData struct {
	Instances []AuthProductInstance `json:"instances"`
}

type AuthProductInstance struct {
	AppMapping map[string]bool `json:"appMapping"`
	PIID       uuid.UUID       `json:"piid"`
}
