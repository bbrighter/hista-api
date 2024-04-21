package meals

import "context"

type IngredientParams struct {
	Name string `json:"name"`
}

type IDResponse struct {
	ID uint `json:"id"`
}

// encore:api public method=PUT path=/ingredient
func (service Service) PutIngredient(ctx context.Context, params IngredientParams) (IDResponse, error) {
	id, err := service.createOrReplaceIngredient(params.Name)
	return IDResponse{ID: id}, err
}

type IngredientResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type IngredientsResponse struct {
	Ingredients []IngredientResponse `json:"ingredients"`
}

func ingredientToIngredientResponse(ingredient Ingredient) IngredientResponse {
	return IngredientResponse{
		ID:   ingredient.ID,
		Name: ingredient.Name,
	}
}

func ingredientsToIngredientsResponse(ingredients []Ingredient) IngredientsResponse {
	var resps []IngredientResponse
	for _, ing := range ingredients {
		resps = append(resps, ingredientToIngredientResponse(ing))
	}
	return IngredientsResponse{resps}
}

// encore:api public method=GET path=/ingredient
func (service Service) GetIngredients(ctx context.Context) (IngredientsResponse, error) {
	ingredients := service.getIngredients()
	return ingredientsToIngredientsResponse(ingredients), nil
}
