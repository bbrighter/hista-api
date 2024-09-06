package internal

import "encore.app/entity"

type IngredientUseCase struct {
	repo IIngredientRepository
}

func NewIngredientUseCase(repo IIngredientRepository) IngredientUseCase {
	return IngredientUseCase{repo: repo}
}

func (uc IngredientUseCase) List() entity.Ingredients {
	return uc.repo.ListIngredients()
}
