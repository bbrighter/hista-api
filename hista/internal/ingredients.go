package internal

import (
	"context"
	"strings"

	"encore.app/errors"
	"encore.app/hista/entity"
)

type (
	IIngredientRepository interface {
		ListIngredients(ctx context.Context) ([]*entity.Ingredient, error)
		DeleteIngredient(ctx context.Context, id uint) error
		UpdateIngredient(ctx context.Context, id uint, values map[string]any) error
		ToggleArchived(ctx context.Context, id uint) error
	}

	IIngredientUseCase interface {
		List(ctx context.Context) (entity.Ingredients, error)
	}

	IIngredientManager interface {
		ChangeIngredient(ctx context.Context, id uint, newName *string, newNutrition *entity.Nutrition, isArchived *bool) error
		ToggleArchived(ctx context.Context, id uint) error
		Delete(ctx context.Context, id uint) error
	}
)

type IngredientUseCase struct {
	repo IIngredientRepository
}

func NewIngredientUseCase(repo IIngredientRepository) IngredientUseCase {
	return IngredientUseCase{repo: repo}
}

func (uc IngredientUseCase) List(ctx context.Context) (entity.Ingredients, error) {
	ings, err := uc.repo.ListIngredients(ctx)
	return ings, errors.MapError(err)
}

type IngredientManager struct {
	repo IIngredientRepository
}

func NewIngredientsManager(repo IIngredientRepository) IngredientManager {
	return IngredientManager{repo: repo}
}

func (uc IngredientManager) ChangeIngredient(ctx context.Context, id uint,
	newName *string,
	newNutrition *entity.Nutrition,
	isArchived *bool,
) error {
	values := make(map[string]any)
	if newName != nil {
		values["name"] = strings.TrimSpace(*newName)
	}
	if newNutrition != nil {
		values["nutrition_protein"] = newNutrition.Protein
		values["nutrition_carbohydrate"] = newNutrition.Carbohydrate
		values["nutrition_fat"] = newNutrition.Fat
		values["nutrition_fiber"] = newNutrition.Fiber
	}
	if isArchived != nil {
		values["is_archived"] = *isArchived
	}
	return uc.repo.UpdateIngredient(ctx, id, values)
}

func (uc IngredientManager) ToggleArchived(ctx context.Context, id uint) error {
	return uc.repo.ToggleArchived(ctx, id)
}

func (uc IngredientManager) Delete(ctx context.Context, id uint) error {
	return uc.repo.DeleteIngredient(ctx, id)
}
