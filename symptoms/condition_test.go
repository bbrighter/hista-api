package symptoms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConditions(t *testing.T) {
	service := initTest(t)

	var conditions Conditions = getConditions(service, 1)

	assert.Len(t, conditions, 1)
	assert.EqualValues(t, 1, conditions[0].Symptom.ID)
}

func TestCreateConditionBySymptomName(t *testing.T) {
	service := initTest(t)

	var err error
	var categories SymptomCategories

	var condition = &Condition{Severity: High, ConditionEventID: 1}

	categories, err = condition.createConditionBySymptomName(service, "Name", 1000)
	assert.Error(t, err)

	condition = &Condition{Severity: High, ConditionEventID: 1}
	categories, err = condition.createConditionBySymptomName(service, "New Name", 1)
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

	var condition = &Condition{SymptomID: 1, ConditionEventID: 100}
	err = condition.createConditionBySymptomID(service)
	assert.Error(t, err)

	condition = &Condition{SymptomID: 1, ConditionEventID: 1}

	err = condition.createConditionBySymptomID(service)
	assert.NoError(t, err)

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

	condition = &Condition{ID: 1}
	err = condition.delete(service)
	assert.NoError(t, err)
}

func TestChangeSeverity(t *testing.T) {
	service := initTest(t)

	var condition = &Condition{ID: 1}
	var err error = condition.changeSeverity(service, Low)
	assert.NoError(t, err)
	assert.Equal(t, Low, condition.Severity)

	service.db.Find(condition)
	assert.Equal(t, Low, condition.Severity)

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
