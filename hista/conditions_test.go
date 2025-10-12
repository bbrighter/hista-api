package hista

import (
	"testing"

	"encore.app/hista/entity"
	"github.com/stretchr/testify/assert"
)

func TestPatchCondition(t *testing.T) {
	t.Skip()
	service, ctx := initAPITest(t)

	id := service.createTestEvent(ctx, t)

	var err error
	err = service.PatchCondition(ctx, TEST_PIID, 100, PatchSeverityRequestParams{Severity: entity.HighSeverity})
	assert.EqualError(t, err, "not_found: not found")

	catId, err := service.symptoms.CreateCategory(ctx, "cat")
	assert.NoError(t, err)
	symtpomName := "name"
	resp, err := service.PostCondition(ctx, TEST_PIID, id, ConditionRequestParams{SymptomName: &symtpomName, CategoryID: &catId})
	defer service.DeleteCondition(ctx, TEST_PIID, resp.Condition.ID)
	assert.NoError(t, err)
	conditionId := resp.Condition.ID

	err = service.PatchCondition(ctx, TEST_PIID, conditionId, PatchSeverityRequestParams{Severity: entity.HighSeverity})
	assert.NoError(t, err)
}

func TestDeleteCondition(t *testing.T) {
	service, ctx := initAPITest(t)

	id := service.createTestEvent(ctx, t)
	catId, err := service.symptoms.CreateCategory(ctx, "cat2")
	assert.NoError(t, err)
	symtpomName := "name"
	resp, err := service.PostCondition(ctx, TEST_PIID, id, ConditionRequestParams{SymptomName: &symtpomName, CategoryID: &catId})
	assert.NoError(t, err)
	conditionId := resp.Condition.ID

	cats, err := service.DeleteCondition(ctx, TEST_PIID, conditionId)
	assert.NoError(t, err)

	assert.Len(t, cats.Categories, 1)
}
