package api

import (
	"context"

	entity "encore.app/entity"
	"encore.app/internal/repositories/meals"
	"encore.dev/types/uuid"
)

// encore:api auth method=DELETE path=/piid/:piid/foods/:foodId tag:external
func (service *Service) DeleteFood(ctx context.Context, piid uuid.UUID, foodId uint) (entity.IngredientsResponse, error) {
	ingredients, err := service.foods.Delete(ctx, foodId)
	return ingredients.ToIngredientsResponse(), err
}

type FoodConditionParams struct {
	Condition string `query:"condition"`
}

// encore:api auth method=PATCH path=/piid/:piid/foods/:foodId/condition tag:external
func (service *Service) PatchFoodCondition(ctx context.Context, piid uuid.UUID, foodId uint, params FoodConditionParams) error {
	var err error
	condition, err := meals.StringToFoodCondition(params.Condition)
	if err != nil {
		return err
	}
	return service.foods.ChangeCondition(ctx, foodId, condition)

}
