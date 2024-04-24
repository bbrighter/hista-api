package meals

import "context"

type FoodResponse struct {
	ID         uint               `json:"id"`
	Ingredient IngredientResponse `json:"ingredient"`
	Condition  FoodCondition      `json:"foodCondition"`
}

type FoodsResponse struct {
	Foods []FoodResponse `json:"foods"`
}

func foodsToFoodsResponse(foods []Food) []FoodResponse {
	var foodsResponse = []FoodResponse{}
	for _, f := range foods {
		foodsResponse = append(foodsResponse,
			FoodResponse{
				ID:         f.ID,
				Ingredient: ingredientToIngredientResponse(f.Ingredient),
				Condition:  f.Condition,
			})
	}
	return foodsResponse
}

// encore:api auth method=GET path=/meal/:mealId/foods
func (service Service) GetFoods(ctx context.Context, mealId uint) (FoodsResponse, error) {
	var foods []Food = service.getFoods(mealId)
	var resp []FoodResponse = foodsToFoodsResponse(foods)
	return FoodsResponse{Foods: resp}, nil
}

type FoodParams struct {
	IngredientName string        `json:"ingredientName"`
	Condition      FoodCondition `json:"condition" validate:"oneof=raw cooked"`
}

// encore:api auth method=POST path=/meal/:mealId/foods
func (service Service) PostFood(ctx context.Context, mealId uint, params FoodParams) (IDResponse, error) {
	food, err := service.createFood(mealId, params.Condition, params.IngredientName)
	return IDResponse{ID: food.ID}, err
}

// encore:api auth method=DELETE path=/meal/:mealId/foods/:foodId
func (service Service) DeleteFood(ctx context.Context, mealId uint, foodId uint) error {
	var food = Food{ID: foodId, MealID: mealId}
	return service.deleteFood(food)
}

type FoodConditionParams struct {
	Condition string `query:"condition"`
}

// encore:api auth method=PATCH path=/meal/:mealId/foods/:foodId/condition
func (service Service) PatchFoodCondition(ctx context.Context, mealId uint, foodId uint, params FoodConditionParams) error {
	var food = Food{ID: foodId}
	condition, err := stringToFoodCondition(params.Condition)
	if err != nil {
		return err
	}
	return service.changeFoodCondition(food, condition)
}
