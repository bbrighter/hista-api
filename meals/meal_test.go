package meals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateMeal(t *testing.T) {
	service := initTest(t)

	var err error

	var meal = Meal{Date: time.Date(1999, 0, 0, 0, 0, 0, 0, time.Local)}
	err = meal.create(service)
	defer meal.delete(service)

	assert.GreaterOrEqual(t, meal.ID, uint(1))
	assert.NoError(t, err)
	var mealInDB Meal
	service.db.First(&mealInDB, Meal{ID: meal.ID})
	assert.True(t, meal.Date.Equal(mealInDB.Date))
}

func TestGetMeals(t *testing.T) {
	service := initTest(t)

	var meals Meals
	var err error
	err = meals.get(service)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(meals), 1)
}

func TestGetMeal(t *testing.T) {
	service := initTest(t)

	var err error
	var meal = Meal{ID: 1000}

	err = meal.get(service)
	assert.Error(t, err)

	meal.ID = testMeal.ID
	err = meal.get(service)
	assert.NoError(t, err)
	assert.Equal(t, testMeal.Freshness, meal.Freshness)
	assert.GreaterOrEqual(t, len(meal.Foods), 1)
}

func TestDeleteMeal(t *testing.T) {
	service := initTest(t)

	var err error
	var nonExistingMeal = Meal{ID: 100}
	err = nonExistingMeal.delete(service)
	assert.Error(t, err)

	var meal Meal
	err = meal.create(service)
	assert.NoError(t, err)
	err = meal.delete(service)
	assert.NoError(t, err)
}

func TestPatchMeal(t *testing.T) {
	service := initTest(t)

	var meal Meal
	var err error
	err = meal.create(service)
	defer meal.delete(service)
	assert.NoError(t, err)

	var params = PatchParams{}
	var patchDate time.Time = time.Date(1700, 0, 0, 0, 0, 0, 0, time.Local)
	var patchFreshness Freshness = Older
	params.Date = &patchDate
	err = meal.patch(service, params)
	assert.NoError(t, err)
	params.Freshness = &patchFreshness
	err = meal.patch(service, params)
	assert.NoError(t, err)
	var patchStressLevel uint8
	params.StressLevel = &patchStressLevel
	err = meal.patch(service, params)
	assert.NoError(t, err)

	var mealInDB = Meal{ID: meal.ID}
	service.db.First(&mealInDB)
	assert.True(t, mealInDB.Date.Equal(patchDate))
	assert.Equal(t, patchFreshness, mealInDB.Freshness)
	assert.Equal(t, patchStressLevel, mealInDB.StressLevel)
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
