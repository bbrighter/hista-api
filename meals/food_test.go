package meals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFood(t *testing.T) {
	service := initTest(t)

	var foods []Food = getFoods(service, testMeal.ID)
	assert.GreaterOrEqual(t, len(foods), 1)
}

func TestCreateFoodByName(t *testing.T) {
	service := initTest(t)

	var err error
	var food = &Food{MealID: testMeal.ID}
	var ingredients Ingredients

	ingredients, err = food.createByName(service, "New name")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(ingredients), 1)

	// Cleanup
	_, err = food.delete(service)
	assert.NoError(t, err)

	// Food for non-existing meal
	var nonExistingFood = &Food{MealID: 10000}
	_, err = nonExistingFood.createByName(service, "New name 2")

	assert.Error(t, err)
}

func TestCreateFoodById(t *testing.T) {
	service := initTest(t)

	var food = &Food{MealID: testMeal.ID, IngredientID: testFood.IngredientID}
	var ingredients Ingredients
	var err error
	ingredients, err = food.createByID(service)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(ingredients), 1)

	// Cleanup
	err = service.db.Delete(&food).Error
	assert.NoError(t, err)
}

func TestDeleteFood(t *testing.T) {
	service := initTest(t)

	var err error
	var nonExistingFood, food Food
	var ingredients Ingredients
	food.ID = testFood.ID
	ingredients, err = food.delete(service)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(ingredients), 0)
	rows := service.db.Find(&Ingredient{ID: food.IngredientID}).RowsAffected
	assert.EqualValues(t, 0, rows)

	// Test error
	nonExistingFood.ID = 100000
	_, err = nonExistingFood.delete(service)
	assert.Error(t, err)
}

func TestChangeFoodCondition(t *testing.T) {
	service := initTest(t)

	var err error
	err = testFood.changeCondition(service, Raw)
	assert.NoError(t, err)

	var nonExistingFood = &Food{ID: 10000}
	err = nonExistingFood.changeCondition(service, Raw)
	assert.Error(t, err)
}

func TestStringToFoodCondition(t *testing.T) {
	t.Parallel()

	var err error
	var cond FoodCondition
	cond, err = stringToFoodCondition("raw")
	assert.NoError(t, err)
	assert.Equal(t, Raw, cond)

	_, err = stringToFoodCondition("bad input")
	assert.Error(t, err)
}
