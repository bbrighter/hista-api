package symptoms

import (
	"context"

	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

type SymptomRepo struct {
	db *gorm.DB
}

func NewSymptomRepo(db *gorm.DB) *SymptomRepo {
	return &SymptomRepo{db: db}
}

func (r *SymptomRepo) ListSymptomCategoriesAndSymptoms(ctx context.Context) (SymptomCategories, error) {
	return gorm.G[SymptomCategory](r.db).
		Scopes(wherePiid(ctx)).
		Preload("Symptoms", nil).
		Find(ctx)
}

func (r *SymptomRepo) CreateSymptomCategory(ctx context.Context, cat *SymptomCategory) error {
	return generic_queries.Create(ctx, r.db, cat)
}

func (r *SymptomRepo) UpdateSymptom(ctx context.Context, id uint, values map[string]any) error {
	rows, err := gorm.G[map[string]any](r.db).
		Table("symptoms").
		Where("id = ?", id).
		Updates(ctx, values)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SymptomRepo) UpdateSymptomCategory(ctx context.Context, id uint, values map[string]any) error {
	return generic_queries.Updates(ctx, r.db, "symptom_categories", id, values)
}

func (r *SymptomRepo) DeleteSymptomCategory(ctx context.Context, id uint) error {
	return generic_queries.Delete[*SymptomCategory](ctx, r.db, id)
}
