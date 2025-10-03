package internal

import (
	"context"

	"encore.app/entity"
)

type (
	IIngredientRepository interface {
		ListIngredients(ctx context.Context) ([]*entity.Ingredient, error)
	}

	IIngredientUseCase interface {
		List(ctx context.Context) (entity.Ingredients, error)
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
	return ings, errorMapper(err)
}
