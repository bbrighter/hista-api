package api

import (
	"context"

	"encore.app/entity"
)

// encore:api auth method=GET path=/ingredients
func (service *Service) GetIngredients(ctx context.Context) (entity.IngredientsResponse, error) {
	return service.ingredients.List().ToIngredientsResponse(), nil
}
