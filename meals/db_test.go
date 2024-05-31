package meals

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testMeal *Meal
var testFood *Food
var testIngredient *Ingredient

func initAPITest(t *testing.T) (*Service, context.Context) {
	var ctx context.Context = context.TODO()
	service := initTest(t)
	return service, ctx
}

func initTest(t *testing.T) *Service {
	service, err := initService()
	assert.NoError(t, err)
	service.initData()
	return service
}

func (service *Service) initData() {
	var meal = &Meal{
		Date:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
		Freshness:   Fresh,
		StressLevel: 3,
		IsAlone:     true,
	}
	service.db.FirstOrCreate(&meal, &meal)
	var ingredient = &Ingredient{Name: "Ingredient"}
	service.db.FirstOrCreate(&ingredient, &ingredient)
	var food = &Food{IngredientID: ingredient.ID, Condition: Cooked, MealID: meal.ID}
	service.db.FirstOrCreate(&food, &food)

	testMeal = meal
	testFood = food
	testIngredient = ingredient
}
