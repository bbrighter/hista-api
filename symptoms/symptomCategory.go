package symptoms

import (
	"encore.app/errors"
	"gorm.io/gorm/clause"
)

type SymptomCategory struct {
	ID       uint
	Name     string
	Symptoms []Symptom
}

type SymptomCategories []SymptomCategory

// Get all symptom categories and its children
func getSymptomCategories(service *Service) SymptomCategories {
	var cats []SymptomCategory
	service.db.Preload(clause.Associations).Find(&cats)
	return cats
}

// Create a new symptom category
// Name is required
func (cat *SymptomCategory) create(service *Service) error {
	if cat.Name == "" {
		return errors.ErrorAttributeMustBeSet("name")
	}

	return service.db.FirstOrCreate(cat, cat).Error
}

// // Update the name of symptom category by ID
// func (cat *SymptomCategory) update(service *Service, newName string) error {
// 	if cat.ID == 0 {
// 		return errors.ErrorIDMissing
// 	}
// 	tx := service.db.Where(cat).Updates(SymptomCategory{Name: newName})
// 	if tx.RowsAffected == 0 {
// 		return errors.ErrorNotFound
// 	}
// 	return tx.Error
// }

// // Delete a symptom category by ID
// func (cat *SymptomCategory) delete(service *Service) error {
// 	if cat.ID == 0 {
// 		return errors.ErrorIDMissing
// 	}
// 	var symptoms Symptoms
// 	if rowsAffected := service.db.
// 		Where(&Symptom{SymptomCategoryID: cat.ID}).
// 		Find(&symptoms).RowsAffected; rowsAffected > 0 {
// 		return errors.NewError("cannot be deleted, symptoms exist", 400)
// 	}
// 	tx := service.db.Delete(cat)
// 	if tx.RowsAffected == 0 {
// 		return errors.ErrorNotFound
// 	}
// 	return tx.Error
// }
