package meals

import "gorm.io/gorm"

type Ingredient struct {
	ID   uint
	Name string `gorm:"uniqueIndex"`
}

type Ingredients []Ingredient

func (service Service) createOrReplaceIngredient(name string) (Ingredient, error) {
	var ingredient = Ingredient{Name: name}
	err := service.db.FirstOrCreate(&ingredient, Ingredient{Name: name}).Error
	return ingredient, err
}

func (service Service) getIngredients() Ingredients {
	var ingredients = []Ingredient{}
	service.db.Find(&ingredients)
	return ingredients
}

func deleteIngredientIfUnused(tx *gorm.DB, ingredientId uint) error {
	var foods []Food
	var err error
	if usedIngredients := tx.Where(Food{IngredientID: ingredientId}).Find(&foods).RowsAffected; usedIngredients == 0 {
		err = tx.Where(&Ingredient{ID: ingredientId}).Delete(&Ingredient{}).Error
	}
	return err
}
