package symptoms

import (
	"encore.app/entity"
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
func (repo *SymptomsRepo) CreateCategory(name string) (uint, error) {
	var cat = entity.SymptomCategory{Name: name}
	err := repo.db.FirstOrCreate(&cat, &cat).Error
	return cat.ID, err
}
