package meals

type Food struct {
	ID           uint
	Ingredient   Ingredient
	IngredientID uint
	Condition    FoodCondition
	MealID       uint
}

type FoodCondition string

const (
	Raw    FoodCondition = "roh"
	Cooked FoodCondition = "gekocht"
)
