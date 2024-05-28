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
	service, ctx := initAPITest(t)

	var params = PatchSeverityRequestParams{Severity: High}
	var err error
	err = service.PatchCondition(ctx, testCondition.ID, params)
	assert.NoError(t, err)

	err = service.PatchCondition(ctx, 1000, params)
	assert.Error(t, err)
}

func TestDeleteConditionAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	var cats SymptomCategoriesResponse

	var condition = Condition{SymptomID: testSymptom.ID, ConditionEventID: testEvent.ID, Severity: VeryHigh}
	service.db.Create(&condition)

	cats, err = service.DeleteCondition(ctx, condition.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(cats.Categories), 1)
	cat1 := cats.Categories[0]
	assert.GreaterOrEqual(t, len(cat1.Symptoms), 1)
}
