package symptoms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConditions(t *testing.T) {
	service := initTest(t)

	var conditions Conditions = getConditions(service, testEvent.ID)

	assert.GreaterOrEqual(t, len(conditions), 1)
	assert.Equal(t, testCondition.ID, conditions[0].Symptom.ID)
}

func TestCreateConditionBySymptomName(t *testing.T) {
	service := initTest(t)

	var err error
	var categories SymptomCategories

	var condition = &Condition{Severity: High, ConditionEventID: testEvent.ID}

	categories, err = condition.createConditionBySymptomName(service, "Name", 1000)
	assert.Error(t, err)

	condition = &Condition{Severity: High, ConditionEventID: testEvent.ID}
	categories, err = condition.createConditionBySymptomName(service, "New Name", testCategory.ID)
	assert.NoError(t, err)
	assert.Equal(t, condition.Symptom.Name, "New Name")
	assert.GreaterOrEqual(t, len(categories), 1)

	// clean up
	err = condition.delete(service)
	assert.NoError(t, err)
}

func TestCreateConditionByID(t *testing.T) {
	service := initTest(t)

	var err error

	var condition = &Condition{SymptomID: 1, ConditionEventID: 10000}
	_, err = condition.createConditionBySymptomID(service)
	assert.Error(t, err)

	condition = &Condition{SymptomID: testSymptom.ID, ConditionEventID: testEvent.ID}

	var symptoms SymptomCategories
	symptoms, err = condition.createConditionBySymptomID(service)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(symptoms), 1)

	// clean up
	err = condition.delete(service)
	assert.NoError(t, err)
}

func TestDeleteCondition(t *testing.T) {
	service := initTest(t)

	var condition = &Condition{ConditionEventID: 1000}
	var err error
	err = condition.delete(service)
	assert.Error(t, err)

	err = testCondition.delete(service)
	assert.NoError(t, err)
}

func TestChangeSeverity(t *testing.T) {
	service := initTest(t)

	var condition = &Condition{ID: testCondition.ID}
	var err error = condition.changeSeverity(service, Low)
	assert.NoError(t, err)
	assert.Equal(t, Low, condition.Severity)

	service.db.Find(condition)
	assert.Equal(t, Low, condition.Severity)

	// cleanup
	condition.changeSeverity(service, High)

	// Error
	condition = &Condition{Severity: High, ConditionEventID: 1000}
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
