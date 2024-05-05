package states

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCondition(t *testing.T) {
	t.Parallel()

	var condition *Condition = newCondition(High, 1)
	assert.Equal(t, High, condition.Severity)
	assert.Equal(t, uint(1), condition.StateID)
	assert.Equal(t, uint(0), condition.ID)
	assert.Equal(t, Symptom{}, condition.Symptom)
}

func TestGetConditions(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var state State = service.testCreateState(t)
	var conditions Conditions = getConditions(service, state.ID)

	assert.Len(t, conditions, 1)
	assert.EqualValues(t, 1, conditions[0].Symptom.ID)
}

func TestCreateCondition(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var err error
	var symptomCategory = SymptomCategory{ID: 2, Name: "Category"}
	err = service.db.Create(&symptomCategory).Error
	assert.NoError(t, err)

	var condition = newCondition(High, 1)

	err = condition.create(service, "Name", 2)
	assert.Error(t, err)

	var state State = service.testCreateState(t)
	condition.StateID = state.ID
	err = condition.create(service, "New Name", 1)
	assert.NoError(t, err)
}

func TestDeleteCondition(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var condition = newCondition(High, 1)
	var err error
	err = condition.delete(service)
	assert.Error(t, err)

	var state State = service.testCreateState(t)
	condition = &state.Conditions[0]
	err = condition.delete(service)
	assert.NoError(t, err)
}

func TestChangeSeverity(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var state State = service.testCreateState(t)
	var condition *Condition = &state.Conditions[0]
	var err error = condition.changeSeverity(service, Low)
	assert.NoError(t, err)
	assert.Equal(t, Low, condition.Severity)

	service.db.Find(condition)
	assert.Equal(t, Low, condition.Severity)

	condition = newCondition(High, 10000)
	err = condition.changeSeverity(service, VeryHigh)
	assert.Error(t, err)
}

func TestNumberToSeverity(t *testing.T) {
	t.Parallel()

	var severity ConditionSeverity
	var err error

	severity, err = numberToSeverity(1)
	assert.NoError(t, err)
	assert.Equal(t, VeryLow, severity)

	_, err = numberToSeverity(100)
	assert.Error(t, err)
}
