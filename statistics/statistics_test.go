package statistics

import (
	"testing"
	"time"

	"encore.app/meals"
	"github.com/stretchr/testify/assert"
)

func TestFindFoodForSymptoms(t *testing.T) {
	service := initTest(t)

	var meal meals.Meal
	service.db.First(&meal)
	var mealTime time.Time = meal.Date

	var fromDate, toDate time.Time
	var symptomIds = []uint{1}
	fromDate = mealTime.Add(-time.Hour * 72)
	toDate = mealTime.Add(time.Hour)
	var results Statistics
	var err error
	results, err = findFoodForSymptoms(service, fromDate, toDate, symptomIds)
	assert.NoError(t, err)
	assert.Len(t, results.Statistics, 1)
	var res Statistic = results.Statistics[0]
	assert.EqualValues(t, 1, res.IngredientID)
	assert.Equal(t, "cooked", res.FoodCondition)
	assert.Len(t, res.Statistic, 1)
	var stat = res.Statistic[0]
	assert.True(t, mealTime.Equal(stat.MealDate), mealTime.String(), stat.MealDate.String())
	// assert.True(t, symptomTime.Equal(stat.ConditionEventDate))
	assert.EqualValues(t, 1, stat.SymptomID)
	assert.Equal(t, 4, stat.SymptomSeverity)

	var futureDate time.Time = mealTime.Add(1000 * time.Hour)
	results, err = findFoodForSymptoms(service, futureDate, futureDate.Add(time.Hour), symptomIds)
	assert.NoError(t, err)
	assert.Len(t, results.Statistics, 0)
}
