package symptoms

import (
	"strings"

	"encore.app/entity"
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
