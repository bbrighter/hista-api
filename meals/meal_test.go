package meals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateMeal(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var id uint
	var err error

	var testIngredient = Ingredient{ID: 100, Name: "Test ingredient"}
	service.db.Create(&testIngredient)
	var date = time.Now()

	id, err = service.createMeal(date)

	assert.GreaterOrEqual(t, id, uint(1))
	assert.NoError(t, err)
}

func TestGetMeals(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var meals []Meal
	var err error
	meals, err = service.getMeals()
	assert.NoError(t, err)
	assert.Len(t, meals, 0)

	service.testCreateMeal(t)

	meals, err = service.getMeals()
	assert.NoError(t, err)
	assert.Len(t, meals, 1)
}

func TestGetMeal(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var err error

	_, err = service.getMeal(1000)
	assert.Error(t, err)

	var meal Meal = service.testCreateMeal(t)
	_, err = service.getMeal(meal.ID)
	assert.NoError(t, err)
}

func TestDeleteMeal(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var err error
	err = service.deleteMeal(1)
	assert.Error(t, err)

	var meal Meal = service.testCreateMeal(t)
	err = service.deleteMeal(meal.ID)
	assert.NoError(t, err)
}

func TestGetMealsAndDependencies(t *testing.T) {
	service, teardown := initTest(t)
	service.testCreateMeal(t)
	defer teardown(t)

	var meals Meals
	var err error
	meals, err = GetMealsAndDependencies(service.db)
	assert.NoError(t, err)
	assert.Len(t, meals, 1)
	var meal Meal = meals[0]
	assert.Len(t, meal.Foods, 1)
	assert.Equal(t, meal.Foods[0].Ingredient.Name, "Name")

}
