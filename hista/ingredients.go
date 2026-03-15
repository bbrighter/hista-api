package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/entity"
	"encore.dev/types/option"
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
	Name      option.Option[string]                      `json:"name"`
	Nutrition option.Option[entity.PatchNutritionParams] `json:"nutrition"`
	Archived  option.Option[bool]                        `json:"archived"`
}

// encore:api auth method=PATCH path=/piid/:piid/ingredients/:id
func (service *Service) PatchIngredient(ctx context.Context, piid uuid.UUID, id uint, params PatchIngredientParams) error {
	name := params.Name.PtrOrNil()
	nutrition := params.Nutrition.PtrOrNil().ToNutrition()
	archived := params.Archived.PtrOrNil()
	if name == nil && nutrition == nil && archived == nil {
		return errors.BadRequest("one property must be set")
	}
	return errors.MapError(service.ingredientsManager.ChangeIngredient(ctx, id, name, nutrition, archived))
}

// Deprecated: use normal PATCH endpoint with archived property
// encore:api auth method=PATCH path=/piid/:piid/ingredients/:id/archive
func (service *Service) ArchiveIngredient(ctx context.Context, piid uuid.UUID, id uint) error {
	return errors.MapError(service.ingredientsManager.ToggleArchived(ctx, id))
}
