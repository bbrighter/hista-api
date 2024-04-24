package meals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFood(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var foods []Food

	// empty response
	foods = service.getFoods(1)
	assert.Len(t, foods, 0)

	// one meal one food
	var meal Meal = service.testCreateMeal(t)
	foods = service.getFoods(meal.ID)
	assert.Len(t, foods, 1)
}

func TestCreateFood(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var food Food
	var err error

	// Food for non-existing meal
	food, err = service.createFood(1, Cooked, "IngredientName")
	assert.Error(t, err)

	// Valid food
	var meal Meal = service.testCreateMeal(t)
	food, err = service.createFood(meal.ID, Cooked, "IngredientName")
	assert.NoError(t, err)
	assert.Equal(t, "IngredientName", food.Ingredient.Name)
}

func TestDeleteFood(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var meal Meal = service.testCreateMeal(t)
	var err error
	err = service.deleteFood(meal.Foods[0])

	assert.NoError(t, err)

	err = service.deleteFood(Food{ID: 100})
	assert.Error(t, err)
}

func TestChangeFoodCondition(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var meal Meal = service.testCreateMeal(t)
	var err error
	err = service.changeFoodCondition(meal.Foods[0], Raw)
	assert.NoError(t, err)

	err = service.changeFoodCondition(Food{ID: 100}, Raw)
	assert.Error(t, err)
}
