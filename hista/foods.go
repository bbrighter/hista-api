package hista

import (
	"context"
	"fmt"

	"encore.app/hista/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=DELETE path=/piid/:piid/foods/:foodId
func (service *Service) DeleteFood(ctx context.Context, piid uuid.UUID, foodId uint) (entity.IngredientsResponse, error) {
	ingredients, err := service.foods.Delete(ctx, foodId)
	return ingredients.ToIngredientsResponse(), err
}

type FoodConditionParams struct {
	Condition entity.FoodCondition `json:"condition"`
}

func (f FoodConditionParams) Validate() error {
	switch f.Condition {
	case entity.Cooked, entity.Raw:
		return nil
	default:
		return fmt.Errorf("invalid condition: %s; valid are %s and %s", f.Condition, entity.Cooked, entity.Raw)
	}
}

// encore:api auth method=PATCH path=/piid/:piid/foods/:foodId/condition
func (service *Service) PatchFoodCondition(ctx context.Context, piid uuid.UUID, foodId uint, params FoodConditionParams) error {
	return service.foods.ChangeCondition(ctx, foodId, params.Condition)
}
