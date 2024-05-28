package meals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateMeal(t *testing.T) {
	service := initTest(t)

	var id uint
	var err error

	var date = time.Now()
	id, err = service.createMeal(date)

	assert.GreaterOrEqual(t, id, uint(1))
	assert.NoError(t, err)

	// cleanup
	err = service.deleteMeal(id)
	assert.NoError(t, err)
}

func TestGetMeals(t *testing.T) {
	service := initTest(t)

	var meals []Meal
	var err error
	meals, err = service.getMeals()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(meals), 1)
}

func TestGetMeal(t *testing.T) {
	service := initTest(t)

	var err error

	_, err = service.getMeal(1000)
	assert.Error(t, err)

	_, err = service.getMeal(testMeal.ID)
	assert.NoError(t, err)
}

func TestDeleteMeal(t *testing.T) {
	service := initTest(t)

	var err error
	err = service.deleteMeal(100)
	assert.Error(t, err)

	var id uint
	id, err = service.createMeal(time.Now())
	assert.NoError(t, err)
	err = service.deleteMeal(id)
	assert.NoError(t, err)
}

func TestGetMealsAndDependencies(t *testing.T) {
	service := initTest(t)

	var meals Meals
	var err error
	meals, err = GetMealsAndDependencies(service.db)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(meals), 1)
	var meal Meal = meals[0]
	assert.GreaterOrEqual(t, len(meal.Foods), 1)

	for _, meal := range meals {
		for _, food := range meal.Foods {
			assert.NotEmpty(t, food.Ingredient.Name)
		}
	}

}
