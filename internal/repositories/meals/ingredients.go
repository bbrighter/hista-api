package meals

import (
	"strings"

	"encore.app/entity"
)

func (repo *MealRepository) CreateOrReplaceIngredient(name string) (entity.Ingredient, error) {
	trimmedName := strings.TrimSpace(name)
	var ingredient = entity.Ingredient{Name: trimmedName}
	err := repo.db.FirstOrCreate(&ingredient, entity.Ingredient{Name: trimmedName}).Error
	return ingredient, err
}

func (repo *MealRepository) ListIngredients() entity.Ingredients {
	var ingredients = entity.Ingredients{}
	repo.db.Find(&ingredients)
	return ingredients
}
