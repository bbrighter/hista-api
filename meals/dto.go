package meals

import "sort"

func (ingredient Ingredient) toIngredientResponse() IngredientResponse {
	return IngredientResponse{
		ID:   ingredient.ID,
		Name: ingredient.Name,
	}
}

func (ingredients Ingredients) toIngredientsResponse() IngredientsResponse {
	var resps []IngredientResponse
	for _, ing := range ingredients {
		resps = append(resps, ing.toIngredientResponse())
	}
	return IngredientsResponse{resps}
}

func (food Food) toFoodResponse() FoodResponse {
	return FoodResponse{
		ID:         food.ID,
		Ingredient: food.Ingredient.toIngredientResponse(),
		Condition:  food.Condition,
	}
}

func (foods Foods) toFoodsResponse() []FoodResponse {
	var foodsResponse = []FoodResponse{}
	for _, f := range foods {
		foodsResponse = append(foodsResponse, f.toFoodResponse())
	}
	return foodsResponse
}

func (meal Meal) toMealMetaResponse() MealMetaResponse {
	return MealMetaResponse{
		ID:   meal.ID,
		Date: meal.Date,
	}
}

func (meal Meal) toMealResponse() MealResponse {
	var foodsResponse []FoodResponse
	for _, f := range meal.Foods {
		foodsResponse = append(foodsResponse,
			FoodResponse{
				ID:         f.ID,
				Ingredient: f.Ingredient.toIngredientResponse(),
				Condition:  f.Condition,
			})
	}
	var resp = MealResponse{
		ID:          meal.ID,
		Date:        meal.Date,
		Foods:       Foods(meal.Foods).toFoodsResponse(),
		Freshness:   meal.Freshness,
		StressLevel: meal.StressLevel,
		IsAlone:     meal.IsAlone,
	}
	return resp
}

func (meals Meals) toMealsResponse() MealsResponse {
	var resps []MealMetaResponse
	for _, m := range meals {
		resps = append(resps, m.toMealMetaResponse())
	}
	sort.Slice(resps, func(i, j int) bool {
		return resps[i].Date.Sub(resps[j].Date) > 0
	})
	return MealsResponse{resps}
}
