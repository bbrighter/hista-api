package meals

import (
	"context"

	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

type ingredientRepo struct {
	db *gorm.DB
}

func newIngredientRepo(db *gorm.DB) *ingredientRepo {
	return &ingredientRepo{db: db}
}

func (r *ingredientRepo) ListIngredients(ctx context.Context) ([]*Ingredient, error) {
	return generic_queries.List[*Ingredient](ctx, r.db)
}

func (r *ingredientRepo) UpdateIngredient(ctx context.Context, id uint, values map[string]any) error {
	return generic_queries.Updates(ctx, r.db, "ingredients", id, values)
}

func (r *ingredientRepo) GetIngredientByName(ctx context.Context, name string) (Ingredient, error) {
	return gorm.G[Ingredient](r.db).
		Scopes(wherePiid(ctx)).
		Where("name = ?", name).
		First(ctx)
}
