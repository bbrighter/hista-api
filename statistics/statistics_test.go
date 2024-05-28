package statistics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindFoodForSymptoms(t *testing.T) {
	service := initTest(t)

	var fromDate, toDate, timeOfMeal time.Time
	timeOfMeal = testMeal.Date
	fromDate = timeOfMeal.Add(-time.Hour * 24 * 7) // one week before meal
	toDate = timeOfMeal                            // until meal

	var symptomIds = []uint{testSymptom.ID}
	var err error
	// Regular test
	var results FoodStatisticsResponse
	results, err = findFoodForSymptoms(service, fromDate, toDate, symptomIds)
	assert.NoError(t, err)
	assert.Len(t, results.Statistics, 1)
	var stat StatisticsByFood = results.Statistics[0]
	// assert.EqualValues(t, 1, stat.IngredientID)
	assert.Equal(t, "cooked", stat.FoodCondition)
	assert.Equal(t, 0, stat.Hours1)
	assert.Equal(t, 1, stat.Hours24)
	assert.Equal(t, 1, stat.Hours72)

	// No data
	var futureDate time.Time = timeOfMeal.Add(time.Hour * 24 * 365)
	results, err = findFoodForSymptoms(service, futureDate, futureDate.Add(time.Hour), symptomIds)
	assert.NoError(t, err)
	assert.Len(t, results.Statistics, 0)

	// // Add more data and run complex test
	// var timeOfNewMeal = timeOfSymptom.Add(-time.Hour * 24 * 2) // 2 days before symptom
	// var newMeal = meals.Meal{Date: timeOfNewMeal}
	// err = service.db.Create(&newMeal).Error
	// assert.NoError(t, err)
	// var food = meals.Food{IngredientID: 1, Condition: meals.Cooked, MealID: newMeal.ID}
	// err = service.db.Create(&food).Error
	// assert.NoError(t, err)

	// results, err = findFoodForSymptoms(service, fromDate, toDate, symptomIds)
	// assert.NoError(t, err)
	// assert.Len(t, results.Statistics, 1)
	// stat = results.Statistics[0]
	// assert.Equal(t, stat.Hours1, 0)
	// assert.Equal(t, stat.Hours24, 1)
	// assert.Equal(t, stat.Hours72, 2)

	// // Cleanup
	// service.db.Delete(&food)
	// service.db.Delete(&newMeal)
}

func TestFindSymptomsForFoods(t *testing.T) {
	service := initTest(t)
	// Get dates from database
	var fromDate, toDate, timeOfSymptom time.Time
	timeOfSymptom = testEvent.Date
	fromDate = timeOfSymptom.Add(-time.Hour * 24 * 7) // one week before meal
	toDate = timeOfSymptom

	var err error
	var ids = []uint{testIngredient.ID}
	// Regular test
	var results SymptomStatisticsResponse
	results, err = findSymptomsForFoods(service, fromDate, toDate, ids)
	assert.NoError(t, err)
	assert.Len(t, results.Statistics, 1)
	var stat StatisticBySymptom = results.Statistics[0]
	assert.EqualValues(t, testSymptom.ID, stat.SymptomID)
	assert.Equal(t, 4, stat.Severity)
	assert.Equal(t, 0, stat.Hours1)
	assert.Equal(t, 1, stat.Hours24)
	assert.Equal(t, 1, stat.Hours72)
}
