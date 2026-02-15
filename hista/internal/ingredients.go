package internal

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/entity"
)

type (
	IIngredientRepository interface {
		ListIngredients(ctx context.Context) ([]*entity.Ingredient, error)
		ChangeIngredientName(ctx context.Context, id uint, newName string) error
		ToggleArchived(ctx context.Context, id uint) error
		DeleteIngredient(ctx context.Context, id uint) error
	}

	IIngredientUseCase interface {
		List(ctx context.Context) (entity.Ingredients, error)
	}

	IIngredientManager interface {
		ChangeName(ctx context.Context, id uint, newName string) error
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

func (uc IngredientManager) ChangeName(ctx context.Context, id uint, newName string) error {
	return uc.repo.ChangeIngredientName(ctx, id, newName)
}
func (uc IngredientManager) ToggleArchived(ctx context.Context, id uint) error {
	return uc.repo.ToggleArchived(ctx, id)
}
func (uc IngredientManager) Delete(ctx context.Context, id uint) error {
	return uc.repo.DeleteIngredient(ctx, id)
}
