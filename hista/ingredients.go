package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/ingredients
func (service *Service) ListIngredients(ctx context.Context, piid uuid.UUID) (entity.IngredientsResponse, error) {
	ing, err := service.ingredients.List(ctx)
	return ing.ToIngredientsResponse(), errors.MapError(err)
}

// encore:api auth method=DELETE path=/piid/:piid/ingredients/:id
func (service *Service) DeleteIngredient(ctx context.Context, piid uuid.UUID, id uint) error {
	return errors.MapError(service.ingredientsManager.Delete(ctx, id))
}

type PatchIngredientParams struct {
	Name string `json:"name"`
}

// encore:api auth method=PATCH path=/piid/:piid/ingredients/:id
func (service *Service) PatchIngredient(ctx context.Context, piid uuid.UUID, id uint, params PatchIngredientParams) error {
	return errors.MapError(service.ingredientsManager.ChangeName(ctx, id, params.Name))
}

// encore:api auth method=PATCH path=/piid/:piid/ingredients/:id/archive
func (service *Service) ArchiveIngredient(ctx context.Context, piid uuid.UUID, id uint) error {
	return errors.MapError(service.ingredientsManager.ToggleArchived(ctx, id))
}
