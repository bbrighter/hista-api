package pollen

import (
	"gorm.io/gorm"
)

type PollenRepo struct {
	db *gorm.DB
}

func NewPollenRepo(db *gorm.DB) *PollenRepo {
	return &PollenRepo{db: db}
}
