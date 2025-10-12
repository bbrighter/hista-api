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
