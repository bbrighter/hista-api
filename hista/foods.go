package hista

import (
	"context"

	"encore.app/errors"
	"encore.dev/types/uuid"
)

// encore:api auth method=DELETE path=/piid/:piid/foods/:foodId
func (service *Service) DeleteFood(ctx context.Context, piid uuid.UUID, foodId uint) (IngredientsResponse, error) {
	err := service.meals.DeleteFood(ctx, foodId)
	if err != nil {
		return IngredientsResponse{}, errors.MapError(err)
	}
	ingredients, err := service.meals.ListIngredients(ctx)
	if err != nil {
		return IngredientsResponse{}, errors.MapError(err)
	}
	return toIngredientsResponse(ingredients), nil
}

type PatchFoodConditionParams struct {
	Condition string `json:"condition"`
}

func (f PatchFoodConditionParams) Validate() error {
	if f.Condition != "cooked" && f.Condition != "raw" {
		return errors.BadRequestf("invalid condition: %s; valid are %s and %s", f.Condition, "cooked", "raw")
	}
	return nil
}

// encore:api auth method=PATCH path=/piid/:piid/foods/:foodId/condition
func (service *Service) PatchFoodCondition(ctx context.Context, piid uuid.UUID, foodId uint, params PatchFoodConditionParams) error {
	err := service.meals.ChangeFoodCondition(ctx, foodId, params.Condition)
	return errors.MapError(err)
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
	err := service.meals.ChangeFoodAmount(ctx, foodId, params.Amount)
	return errors.MapError(err)
}
