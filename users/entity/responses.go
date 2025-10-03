package entity

import "encore.dev/types/uuid"

type UserResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
}
