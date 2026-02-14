package meals

import (
	"context"
	"strings"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

func (repo *MealRepository) CreateOrReplaceIngredient(ctx context.Context, name string) (entity.Ingredient, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return entity.Ingredient{}, err
	}
	tx := gorm.G[entity.Ingredient](repo.db).Where("pi_id = ?", piid)

	trimmedName := strings.TrimSpace(name)
	ingredient, err := tx.Where("name = ?", trimmedName).First(ctx)
	if err == nil {
		return ingredient, err
	}
	ingredient = entity.Ingredient{Name: trimmedName, PIID: piid}
	err = gorm.G[entity.Ingredient](repo.db).Create(ctx, &ingredient)
	return ingredient, err
}

func (repo *MealRepository) ListIngredients(ctx context.Context) ([]*entity.Ingredient, error) {
	return generic_queries.List[*entity.Ingredient](ctx, repo.db)
}

func (repo *MealRepository) ChangeIngredientName(ctx context.Context, id uint, newName string) error {
	return generic_queries.UpdateColumn[*entity.Ingredient](ctx, repo.db, id, "name", newName)
}
func (repo *MealRepository) ToggleArchived(ctx context.Context, id uint) error {
	return generic_queries.UpdateColumn[*entity.Ingredient](ctx, repo.db, id, "is_archived", gorm.Expr("NOT is_archived"))
}
func (repo *MealRepository) DeleteIngredient(ctx context.Context, id uint) error {
	return generic_queries.Delete[*entity.Ingredient](ctx, repo.db, id)
}
