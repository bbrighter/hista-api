package api

import (
	"testing"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
)

func TestPatchCondition(t *testing.T) {
	t.Skip()
	service, ctx := initAPITest(t)

	cleanup := service.createTestEvent(t)
	defer cleanup(t)

	var err error
	err = service.PatchCondition(ctx, 100, PatchSeverityRequestParams{Severity: entity.High})
	assert.EqualError(t, err, "not_found: not found")

	catId, err := service.symtpoms.CreateCategory("cat")
	assert.NoError(t, err)
	symtpomName := "name"
	resp, err := service.PostCondition(ctx, testEvent.ID, ConditionRequestParams{SymptomName: &symtpomName, CategoryID: &catId})
	defer service.DeleteCondition(ctx, resp.Condition.ID)
	assert.NoError(t, err)
	conditionId := resp.Condition.ID

	err = service.PatchCondition(ctx, conditionId, PatchSeverityRequestParams{Severity: entity.High})
	assert.NoError(t, err)
}

func TestDeleteCondition(t *testing.T) {
	service, ctx := initAPITest(t)

	cleanup := service.createTestEvent(t)
	defer cleanup(t)
	catId, err := service.symtpoms.CreateCategory("cat2")
	assert.NoError(t, err)
	symtpomName := "name"
	resp, err := service.PostCondition(ctx, testEvent.ID, ConditionRequestParams{SymptomName: &symtpomName, CategoryID: &catId})
	assert.NoError(t, err)
	conditionId := resp.Condition.ID

	cats, err := service.DeleteCondition(ctx, conditionId)
	assert.NoError(t, err)

	assert.Len(t, cats.Categories, 0)
}
