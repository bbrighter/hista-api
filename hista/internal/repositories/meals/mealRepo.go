package meals

import (
	"gorm.io/gorm"
)

type MealRepository struct {
	db *gorm.DB
}

func NewMealRepository(db *gorm.DB) *MealRepository {
	return &MealRepository{db: db}
}
