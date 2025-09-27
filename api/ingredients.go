package api

import (
	"context"

	"encore.app/entity"
)

// encore:api auth method=GET path=/ingredients
func (service *Service) GetIngredients(ctx context.Context) (entity.IngredientsResponse, error) {
	ing, err := service.ingredients.List(ctx)
	return ing.ToIngredientsResponse(), err
}
