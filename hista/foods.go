package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=DELETE path=/piid/:piid/foods/:foodId
func (service *Service) DeleteFood(ctx context.Context, piid uuid.UUID, foodId uint) (entity.IngredientsResponse, error) {
	ingredients, err := service.foods.Delete(ctx, foodId)
	return ingredients.ToIngredientsResponse(), errors.MapError(err)
}

type PatchFoodConditionParams struct {
	Condition entity.FoodCondition `json:"condition"`
}

func (f PatchFoodConditionParams) Validate() error {
	if f.Condition != entity.Cooked && f.Condition != entity.Raw {
		return errors.BadRequestf("invalid condition: %s; valid are %s and %s", f.Condition, entity.Cooked, entity.Raw)
	}
	return nil
}

// encore:api auth method=PATCH path=/piid/:piid/foods/:foodId/condition
func (service *Service) PatchFoodCondition(ctx context.Context, piid uuid.UUID, foodId uint, params PatchFoodConditionParams) error {
	return errors.MapError(service.foods.ChangeCondition(ctx, foodId, params.Condition))
}

type PatchFoodAmountParams struct {
	Amount *int `json:"amount"`
}

func (f PatchFoodAmountParams) Validate() error {
	if f.Amount != nil && *f.Amount < 0 {
		return errors.BadRequestf("amount must be positive, not %d", *f.Amount)
	}
	return nil
}

// encore:api auth method=PATCH path=/piid/:piid/foods/:foodId/amount
func (service *Service) PatchFoodAmount(ctx context.Context, piid uuid.UUID, foodId uint, params PatchFoodAmountParams) error {
	return errors.MapError(service.foods.ChangeAmount(ctx, foodId, params.Amount))
}
