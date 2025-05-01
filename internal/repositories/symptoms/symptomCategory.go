package symptoms

import (
	"encore.app/entity"
	"encore.app/errors"
	"encore.dev/beta/errs"
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

func (repo *SymptomsRepo) RenameCategory(cat *entity.SymptomCategory, newName string) error {
	if rowsAffected := repo.db.First(cat).RowsAffected; rowsAffected == 0 {
		return errors.ErrorNotFound
	}
	cat.Name = newName
	return repo.db.Save(cat).Error
}

func (repo *SymptomsRepo) DeleteCategory(catId uint) error {
	var cat entity.SymptomCategory
	if rowsAffects := repo.db.Debug().Preload(clause.Associations).Take(&cat, catId).RowsAffected; rowsAffects == 0 {
		return errors.ErrorNotFound
	}
	if len(cat.Symptoms) > 0 {
		return errors.NewError("cannot delete category with symptoms", errs.FailedPrecondition)
	}

	return repo.db.Delete(&entity.SymptomCategory{}, catId).Error
}
