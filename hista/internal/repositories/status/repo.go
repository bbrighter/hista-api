package status

import "gorm.io/gorm"

type StatusRepo struct {
	db *gorm.DB
}

func NewStatusRepo(db *gorm.DB) *StatusRepo {
	return &StatusRepo{db: db}
}
