package symptoms

import (
	"context"
	"strings"

	"encore.app/entity"
	"encore.app/generic_queries"
	"gorm.io/gorm"
)

// Creates a symptom or replaces it. Equality checked by name and categoryId
func (repo *SymptomsRepo) CreateOrReplace(ctx context.Context, symptomName string, symptomCategoryId uint) (uint, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	trimmedName := strings.TrimSpace(symptomName)
	symptom, err := gorm.G[entity.Symptom](repo.db).
		Where("pi_id = ?", piid).
		Where("symptom_category_id = ?", symptomCategoryId).
		Where("name = ?", trimmedName).
		First(ctx)
	if err == nil {
		return symptom.ID, err
	}
	newSymptom := &entity.Symptom{Name: trimmedName, SymptomCategoryID: symptomCategoryId}
	err = generic_queries.Create(ctx, repo.db, newSymptom)
	return newSymptom.ID, err
}

func (repo *SymptomsRepo) ChangeCategory(ctx context.Context, symptomId, newCategoryId uint) error {
	count, err := generic_queries.Count[*entity.SymptomCategory](ctx, repo.db, newCategoryId)
	if err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return generic_queries.UpdateColumn[*entity.Symptom](ctx, repo.db, symptomId, "symptom_category_id", newCategoryId)
}

func (repo *SymptomsRepo) RenameSymptom(ctx context.Context, symptomId uint, newName string) error {
	return generic_queries.UpdateColumn[*entity.Symptom](ctx, repo.db, symptomId, "name", newName)
}
