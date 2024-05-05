package symptoms

import (
	"encore.app/errors"
	"gorm.io/gorm"
)

type Symptom struct {
	ID                uint
	Name              string
	SymptomCategoryID uint
}

type Symptoms []Symptom

// Creates a symptom or replaces it. Equality checked by name and categoryId
func (symptom *Symptom) createOrReplace(service *Service) error {
	if symptom.Name == "" {
		return errors.ErrorAttributeMustBeSet("name")
	}
	if symptom.SymptomCategoryID == 0 {
		return errors.ErrorAttributeMustBeSet("symptomCategoryId")
	}
	return service.db.FirstOrCreate(&symptom, Symptom{Name: symptom.Name, SymptomCategoryID: symptom.SymptomCategoryID}).Error
}

func (symptom *Symptom) deleteIfUnused(tx *gorm.DB) error {
	var conditions []Condition
	var err error
	if usedConditions := tx.Where(Condition{SymptomID: symptom.ID}).Find(&conditions).RowsAffected; usedConditions == 0 {
		err = tx.Where(&Symptom{ID: symptom.ID}).Delete(&Symptom{}).Error
	}
	return err
}
