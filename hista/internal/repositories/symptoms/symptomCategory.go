package symptoms

import (
	"context"
	"errors"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

// Get all symptom categories and its children
func (repo *SymptomsRepo) ListCategories(ctx context.Context) ([]*entity.SymptomCategory, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return entity.SymptomCategories{}, err
	}
	return gorm.G[*entity.SymptomCategory](repo.db).Preload("Symptoms", nil).Where("pi_id = ?", piid).Find(ctx)
}

// Create a new symptom category
// Name is required
func (repo *SymptomsRepo) CreateCategory(ctx context.Context, catName string) (uint, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	cat, err := gorm.G[entity.SymptomCategory](repo.db).
		Where("pi_id = ?", piid).Where("name = ?", catName).First(ctx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cat := entity.SymptomCategory{Name: catName}
		err = generic_queries.Create(ctx, repo.db, &cat)
		return cat.ID, err
	}
	return cat.ID, err
}

func (repo *SymptomsRepo) RenameCategory(ctx context.Context, cat *entity.SymptomCategory, newName string) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	if rowsAffected := repo.db.Where("pi_id = ?", piid).First(cat).RowsAffected; rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	cat.Name = newName
	return repo.db.Save(cat).Error
}

func (repo *SymptomsRepo) DeleteCategory(ctx context.Context, catId uint) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	count, err := gorm.G[entity.Symptom](repo.db).
		Where("pi_id = ?", piid).
		Where("symptom_category_id = ?", catId).
		Count(ctx, "*")
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete category with symptoms")
	}
	return generic_queries.Delete[*entity.SymptomCategory](ctx, repo.db, catId)
}
