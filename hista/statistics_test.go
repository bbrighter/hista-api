package hista

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStatisticsBySymptomIds(t *testing.T) {
	service, ctx := initAPITest(t)

	from := time.Now().Add(-time.Hour)
	to := time.Now().Add(time.Hour)

	ids := []uint{1}
	_, err := service.GetStatisticsBySymptomIds(ctx, TEST_PIID, StatisticParams{IDs: ids, FromDate: from, ToDate: to})
	assert.NoError(t, err)

	service.createTestFood(ctx, t)
	_, symptomId, _ := service.createTestSymptom(ctx, t)
	ids = append(ids, symptomId)

	resp, err := service.GetStatisticsBySymptomIds(ctx, TEST_PIID, StatisticParams{IDs: ids, FromDate: from, ToDate: to})
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
	_, err := service.GetStatisticsByIngredientsIds(ctx, TEST_PIID, StatisticParams{IDs: ids, FromDate: from, ToDate: to})
	assert.NoError(t, err)

	_, ingredientId := service.createTestFood(ctx, t)
	service.createTestSymptom(ctx, t)
	ids = append(ids, ingredientId)

	resp, err := service.GetStatisticsByIngredientsIds(ctx, TEST_PIID, StatisticParams{IDs: ids, FromDate: from, ToDate: to})
	assert.NoError(t, err)
	require.Len(t, resp.Statistics, 1)
	stat := resp.Statistics[0]
	assert.EqualValues(t, 1, stat.Count)
	assert.EqualValues(t, 1, stat.Hours1)
	assert.EqualValues(t, 1, stat.Hours24)
	assert.EqualValues(t, 1, stat.Hours72)
}
