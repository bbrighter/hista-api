package meals

import "context"

// type IngredientParams struct {
// 	Name string `json:"name"`
// }

type IDResponse struct {
	ID uint `json:"id"`
}

// // encore:api public method=PUT path=/ingredient
// func (service Service) PutIngredient(ctx context.Context, params IngredientParams) (IDResponse, error) {
// 	ingredient, err := service.createOrReplaceIngredient(params.Name)
// 	return IDResponse{ID: ingredient.ID}, err
// }

type IngredientResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type IngredientsResponse struct {
	Ingredients []IngredientResponse `json:"ingredients"`
}

// encore:api auth method=GET path=/ingredient
func (service Service) GetIngredients(ctx context.Context) (IngredientsResponse, error) {
	ingredients := service.getIngredients()
	return ingredients.toIngredientsResponse(), nil
}
