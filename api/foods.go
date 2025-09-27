package api

import (
	"context"

	entity "encore.app/entity"
	"encore.app/internal/repositories/meals"
)

// encore:api auth method=DELETE path=/foods/:foodId
func (service *Service) DeleteFood(ctx context.Context, foodId uint) (entity.IngredientsResponse, error) {
	ingredients, err := service.foods.Delete(ctx, foodId)
	return ingredients.ToIngredientsResponse(), err
}

type FoodConditionParams struct {
	Condition string `query:"condition"`
}

// encore:api auth method=PATCH path=/foods/:foodId/condition
func (service *Service) PatchFoodCondition(ctx context.Context, foodId uint, params FoodConditionParams) error {
	var err error
	condition, err := meals.StringToFoodCondition(params.Condition)
	if err != nil {
		return err
	}
	return service.foods.ChangeCondition(ctx, foodId, condition)

}
