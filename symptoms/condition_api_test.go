package symptoms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestPostCondition(t *testing.T) {
// 	service, ctx, teardown := initAPITest(t)
// 	service.testCreateConditionEvent(t)
// 	defer teardown(t)

// 	var params = ConditionRequestParams{
// 		SymptomName: "new name",
// 		CategoryID:  1,
// 	}
// 	var eventId uint = 1
// 	_, err := service.PostCondition(ctx, eventId, params)
// 	assert.NoError(t, err)

// }

// func TestPostConditionBySymptomID(t *testing.T) {
// 	service, ctx, teardown := initAPITest(t)
// 	var event ConditionEvent = service.testCreateConditionEvent(t)
// 	defer teardown(t)

// 	var eventID uint = event.ID
// 	var symptomID uint = event.Conditions[0].SymptomID

// 	_, err := service.PostConditionBySymptomID(ctx, eventID, symptomID)
// 	assert.NoError(t, err)
// }

func TestPatchConditionAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	var event ConditionEvent = service.testCreateConditionEvent(t)
	defer teardown(t)

	var params = PatchSeverityRequestParams{Severity: High}
	var conditionID uint = event.Conditions[0].ID
	var err error
	err = service.PatchCondition(ctx, conditionID, params)
	assert.NoError(t, err)

	err = service.PatchCondition(ctx, 1000, params)
	assert.Error(t, err)
}

func TestDeleteConditionAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	var event ConditionEvent = service.testCreateConditionEvent(t)
	defer teardown(t)

	var conditionID uint = event.Conditions[0].ID
	var err error
	err = service.DeleteCondition(ctx, conditionID)
	assert.NoError(t, err)
}
