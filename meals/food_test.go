package meals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFood(t *testing.T) {
	service := initTest(t)

	var foods []Food

	foods = service.getFoods(1)
	assert.Len(t, foods, 1)
}

func TestCreateFood(t *testing.T) {
	service := initTest(t)

	var err error

	// Food for non-existing meal
	_, err = service.createFood(100, Cooked, "Meal doesn't exist")
	assert.Error(t, err)

	// Valid food
	var food Food
	food, err = service.createFood(1, Cooked, "New name")
	assert.NoError(t, err)
	assert.Equal(t, "New name", food.Ingredient.Name)

	// Cleanup
	service.deleteFood(food)
}

func TestDeleteFood(t *testing.T) {
	service := initTest(t)

	var err error
	err = service.deleteFood(Food{ID: 1})
	assert.NoError(t, err)

	err = service.deleteFood(Food{ID: 100})
	assert.Error(t, err)
}

func TestChangeFoodCondition(t *testing.T) {
	service := initTest(t)

	var err error
	err = service.changeFoodCondition(Food{ID: 1}, Raw)
	assert.NoError(t, err)

	err = service.changeFoodCondition(Food{ID: 100}, Raw)
	assert.Error(t, err)
}

func TestStringToFoodCondition(t *testing.T) {
	t.Parallel()

	var err error
	var cond FoodCondition
	cond, err = stringToFoodCondition("raw")
	assert.NoError(t, err)
	assert.Equal(t, Raw, cond)

	cond, err = stringToFoodCondition("bad input")
	assert.Error(t, err)
}
