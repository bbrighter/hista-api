package symptoms

import (
	"strings"

	"encore.app/entity"
	"encore.app/errors"
)

// Creates a symptom or replaces it. Equality checked by name and categoryId
func (repo *SymptomsRepo) CreateOrReplace(symptomName string, symptomCategoryId uint) (uint, error) {
	var symptom entity.Symptom
	err := repo.db.
		FirstOrCreate(&symptom,
			entity.Symptom{
				Name:              strings.TrimSpace(symptomName),
				SymptomCategoryID: symptomCategoryId}).
		Error
	return symptom.ID, err
}

func (repo *SymptomsRepo) ChangeCategory(symptomId, newCategoryId uint) error {
	if rowsAffected := repo.db.Take(&entity.SymptomCategory{}, newCategoryId).RowsAffected; rowsAffected == 0 {
		return errors.ErrorNotFound
	}
	tx := repo.db.Model(&entity.Symptom{}).
		Where("id = ?", symptomId).
		UpdateColumn("symptom_category_id", newCategoryId)
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}

func (repo *SymptomsRepo) RenameSymptom(symptomId uint, newName string) error {
	tx := repo.db.Model(&entity.Symptom{}).
		Where("id = ?", symptomId).
		UpdateColumn("name", newName)
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}
