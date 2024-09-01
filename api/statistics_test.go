package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetStatisticsBySymptomIds(t *testing.T) {
	service, ctx := initAPITest(t)

	from := time.Now().Add(-time.Hour)
	to := time.Now().Add(time.Hour)

	ids := []uint{1}
	_, err := service.GetStatisticsBySymptomIds(ctx, StatisticParams{IDs: ids, FromDate: from, ToDate: to})
	assert.NoError(t, err)

	cleanupFood := service.createTestFood(t)
	defer cleanupFood(t)
	cleanup := service.createTestSymptom(t)
	defer cleanup(t)
	ids = append(ids, testSymptom.ID)

	resp, err := service.GetStatisticsBySymptomIds(ctx, StatisticParams{IDs: ids, FromDate: from, ToDate: to})
	assert.NoError(t, err)
	assert.Len(t, resp.Statistics, 1)
	stat := resp.Statistics[0]
	assert.EqualValues(t, 1, stat.Count)
	assert.EqualValues(t, 1, stat.Hours1)
	assert.EqualValues(t, 1, stat.Hours24)
	assert.EqualValues(t, 1, stat.Hours72)
}

func TestGetStatisticsByIngredientsIds(t *testing.T) {
	service, ctx := initAPITest(t)

	from := time.Now().Add(-time.Hour)
	to := time.Now().Add(time.Hour)

	ids := []uint{1}
	_, err := service.GetStatisticsByIngredientsIds(ctx, StatisticParams{IDs: ids, FromDate: from, ToDate: to})
	assert.NoError(t, err)

	cleanupFood := service.createTestFood(t)
	defer cleanupFood(t)
	cleanup := service.createTestSymptom(t)
	defer cleanup(t)
	ids = append(ids, testFood.IngredientID)

	resp, err := service.GetStatisticsByIngredientsIds(ctx, StatisticParams{IDs: ids, FromDate: from, ToDate: to})
	assert.NoError(t, err)
	assert.Len(t, resp.Statistics, 1)
	stat := resp.Statistics[0]
	assert.EqualValues(t, 1, stat.Count)
	assert.EqualValues(t, 1, stat.Hours1)
	assert.EqualValues(t, 1, stat.Hours24)
	assert.EqualValues(t, 1, stat.Hours72)
}
