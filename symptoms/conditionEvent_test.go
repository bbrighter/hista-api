package symptoms

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConditionEvent(t *testing.T) {
	var event *ConditionEvent = newConditionEvent(time.Now())

	assert.GreaterOrEqual(t, time.Now(), event.Date)
}

func TestCreateConditionEvent(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var event *ConditionEvent
	var err error

	err = event.create(service)
	assert.Error(t, err)

	event = newConditionEvent(time.Now())
	err = event.create(service)
	assert.NoError(t, err)
}

func TestGetConditionEvents(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var events ConditionEvents
	events = getConditionEvents(service)
	assert.Len(t, events, 0)

	var event ConditionEvent = service.testCreateConditionEvent(t)
	events = getConditionEvents(service)
	assert.Len(t, events, 1)
	assert.Equal(t, event.ID, events[0].ID)
	assert.False(t, events[0].Date.IsZero(), "date is set; no comparison because time.Now is used")
}

func TestGetConditionEvent(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	service.testCreateConditionEvent(t)
	var event ConditionEvent
	var err error
	event, err = getConditionEvent(service, 1)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, event.ID)

	// Not found
	_, err = getConditionEvent(service, 1000)
	assert.Error(t, err)
}

func TestDeleteConditionEvent(t *testing.T) {
	service, teardown := initTest(t)
	var event ConditionEvent = service.testCreateConditionEvent(t)
	defer teardown(t)

	var err error
	err = event.delete(service)
	assert.NoError(t, err)
	rows := service.db.Find(&ConditionEvents{}).RowsAffected
	assert.EqualValues(t, 0, rows)

	event.ID = 1000
	err = event.delete(service)
	assert.Error(t, err)
}

func TestPatchConditionEvent(t *testing.T) {
	service, teardown := initTest(t)
	var event ConditionEvent = service.testCreateConditionEvent(t)
	defer teardown(t)

	var err error
	var setTime time.Time = time.Date(2000, 1, 1, 1, 1, 1, 0, time.UTC)
	err = event.patch(service, setTime)
	assert.NoError(t, err)
	var result = ConditionEvent{ID: event.ID}
	service.db.Find(&result)
	assert.True(t, setTime.Equal(result.Date))

	var nonexistingEvent = ConditionEvent{ID: 1000}
	err = nonexistingEvent.patch(service, time.Now())
	assert.Error(t, err)
}

func TestGetConditionEventsAndDependencies(t *testing.T) {
	service, teardown := initTest(t)
	service.testCreateConditionEvent(t)
	defer teardown(t)

	var events ConditionEvents
	var cats SymptomCategories
	events, cats = GetConditionEventsAndDependencies(service.db)
	assert.Len(t, events, 1)
	var event ConditionEvent = events[0]
	assert.Len(t, event.Conditions, 1)
	assert.Equal(t, event.Conditions[0].Symptom.Name, "Name")

	assert.Len(t, cats, 1)
}
