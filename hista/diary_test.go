package hista

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDiary(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, err := service.GetDiary(ctx, TEST_PIID)
	assert.NoError(t, err)
	assert.Len(t, resp.Diaries, 0)

	eventResp, err := service.CreateConditionEvent(ctx, TEST_PIID)
	require.NoError(t, err)
	catResp, err := service.PostSymptomCategory(ctx, TEST_PIID, PostSymptomCategoryRequest{Name: "cat"})
	require.NoError(t, err)
	var symptomName = "symptomName"
	_, err = service.PostCondition(ctx, TEST_PIID, eventResp.ID, ConditionRequestParams{SymptomName: &symptomName, CategoryID: &catResp.ID})
	require.NoError(t, err)

	// cleanupFood := service.createTestFood(ctx, t)
	// defer cleanupFood(t)

	// cleanupNote := service.createTestNote(ctx, t)
	// defer cleanupNote(t)

	resp, err = service.GetDiary(ctx, TEST_PIID)
	assert.NoError(t, err)
	assert.Len(t, resp.Diaries, 1)
}
