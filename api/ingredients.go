package api

import (
	"context"

	"encore.app/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/ingredients tag:external
func (service *Service) ListIngredients(ctx context.Context, piid uuid.UUID) (entity.IngredientsResponse, error) {
	ing, err := service.ingredients.List(ctx)
	return ing.ToIngredientsResponse(), err
}
