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

	trimmedName := strings.TrimSpace(name)
	ingredient, err := gorm.G[entity.Ingredient](repo.db).
		Where("pi_id = ?", piid).
		Where("name = ?", trimmedName).
		First(ctx)

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

func (repo *MealRepository) ToggleArchived(ctx context.Context, id uint) error {
	return generic_queries.UpdateColumn[*entity.Ingredient](ctx, repo.db, id, "is_archived", gorm.Expr("NOT is_archived"))
}
func (repo *MealRepository) DeleteIngredient(ctx context.Context, id uint) error {
	return generic_queries.Delete[*entity.Ingredient](ctx, repo.db, id)
}

func (repo *MealRepository) UpdateIngredient(ctx context.Context, id uint, values map[string]any) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	tx := repo.db.Model(&entity.Ingredient{}).Where("id = ?", id).Where("pi_id = ?", piid).Updates(values)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
