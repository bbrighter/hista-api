package hista

import (
	"context"

	"encore.app/hista/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/ingredients
func (service *Service) ListIngredients(ctx context.Context, piid uuid.UUID) (entity.IngredientsResponse, error) {
	ing, err := service.ingredients.List(ctx)
	return ing.ToIngredientsResponse(), err
}
