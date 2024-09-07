package symptoms

import (
	"encore.app/entity"
	"encore.app/errors"
	"gorm.io/gorm/clause"
)

// Get all symptom categories and its children
func (repo *SymptomsRepo) ListCategories() entity.SymptomCategories {
	var cats []entity.SymptomCategory
	repo.db.Preload(clause.Associations).Find(&cats)
	return cats
}

// Create a new symptom category
// Name is required
func (repo *SymptomsRepo) CreateCategory(cat *entity.SymptomCategory) error {
	if cat.Name == "" {
		return errors.ErrorAttributeMustBeSet("Name")
	}
	return repo.db.FirstOrCreate(&cat, &cat).Error
}
