package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/internal/meals"
	"encore.dev/types/option"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/ingredients
func (service *Service) ListIngredients(ctx context.Context, piid uuid.UUID) (IngredientListResponse, error) {
	ings, err := service.meals.ListIngredients(ctx)
	if err != nil {
		return IngredientListResponse{}, errors.MapError(err)
	}
	return toIngredientsResponse(ings), nil
}

type PatchNutritionParams struct {
	Protein      float32 `json:"protein"`
	Carbohydrate float32 `json:"carbohydrate"`
	Fat          float32 `json:"fat"`
	Fiber        float32 `json:"fiber"`
}

type PatchIngredientParams struct {
	Name      option.Option[string]               `json:"name"`
	Nutrition option.Option[PatchNutritionParams] `json:"nutrition"`
	Archived  option.Option[bool]                 `json:"archived"`
}

func (p *PatchNutritionParams) ToNutrition() *meals.Nutrition {
	if p == nil {
		return nil
	}
	return &meals.Nutrition{
		Protein:      &p.Protein,
		Carbohydrate: &p.Carbohydrate,
		Fat:          &p.Fat,
		Fiber:        &p.Fiber,
	}
}

// encore:api auth method=PATCH path=/piid/:piid/ingredients/:id
func (service *Service) PatchIngredient(ctx context.Context, piid uuid.UUID, id uint, params PatchIngredientParams) error {
	name := params.Name.PtrOrNil()
	nutrition := params.Nutrition.PtrOrNil().ToNutrition()
	archived := params.Archived.PtrOrNil()
	if name == nil && nutrition == nil && archived == nil {
		return errors.BadRequest("one property must be set")
	}
	err := service.meals.UpdateIngredient(ctx, id, name, nutrition, archived)
	return errors.MapError(err)
}

// Deprecated: use normal PATCH endpoint with archived property
// encore:api auth method=PATCH path=/piid/:piid/ingredients/:id/archive
func (service *Service) ArchiveIngredient(ctx context.Context, piid uuid.UUID, id uint) error {
	err := service.meals.ToggleIngredientArchived(ctx, id)
	return errors.MapError(err)
}
