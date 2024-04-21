package meals

type Ingredient struct {
	ID   uint
	Name string
}

func (service Service) createOrReplaceIngredient(name string) (uint, error) {
	var ingredient = Ingredient{Name: name}
	err := service.db.FirstOrCreate(&ingredient, Ingredient{Name: name}).Error
	return ingredient.ID, err
}

func (service Service) getIngredients() []Ingredient {
	var ingredients []Ingredient
	service.db.Find(&ingredients)
	return ingredients
}
