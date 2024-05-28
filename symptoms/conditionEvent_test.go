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
	service := initTest(t)

	var event *ConditionEvent
	var err error

	err = event.create(service)
	assert.Error(t, err)

	event = &ConditionEvent{Date: time.Now()}
	err = event.create(service)
	assert.NoError(t, err)

	// clean up
	err = event.delete(service)
	assert.NoError(t, err)
}

func TestGetConditionEvents(t *testing.T) {
	service := initTest(t)

	var events ConditionEvents
	var err error
	events, err = getConditionEvents(service)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(events), 1)
	var eventIDexist bool = false
	for _, event := range events {
		if event.ID == testEvent.ID {
			eventIDexist = true
		}
	}
	assert.True(t, eventIDexist)
	// assert.Equal(t, testEvent.ID, events[0].ID)
	// assert.True(t, testEvent.Date.Equal(events[0].Date))
}

func TestGetConditionEvent(t *testing.T) {
	service := initTest(t)

	var event ConditionEvent
	var err error
	event, err = getConditionEvent(service, testEvent.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, testEvent.ID, event.ID)

	// Not found
	_, err = getConditionEvent(service, 1000)
	assert.Error(t, err)
}

func TestDeleteConditionEvent(t *testing.T) {
	service := initTest(t)

	var err error
	var event = ConditionEvent{ID: testEvent.ID}
	err = event.delete(service)
	assert.NoError(t, err)
	rows := service.db.Find(&ConditionEvent{ID: testEvent.ID}).RowsAffected
	assert.EqualValues(t, 0, rows)

	event.ID = 1000
	err = event.delete(service)
	assert.Error(t, err)
}

func TestPatchConditionEvent(t *testing.T) {
	service := initTest(t)

	var err error
	var setTime time.Time = time.Date(2000, 1, 1, 1, 1, 1, 0, time.UTC)
	var event = ConditionEvent{ID: testEvent.ID}
	err = event.patch(service, setTime)
	assert.NoError(t, err)
	var result = ConditionEvent{ID: event.ID}
	service.db.Find(&result)
	assert.True(t, setTime.Equal(result.Date))

	var nonexistingEvent = ConditionEvent{ID: 1000}
	err = nonexistingEvent.patch(service, time.Now())
	assert.Error(t, err)

	// cleanup
	testEvent.delete(service)
}

func TestGetConditionEventsAndDependencies(t *testing.T) {
	service := initTest(t)

	var events ConditionEvents
	var cats SymptomCategories
	var err error
	events, cats, err = GetConditionEventsAndDependencies(service.db)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(events), 1)
	var event ConditionEvent = events[0]
	assert.GreaterOrEqual(t, len(event.Conditions), 1)

	var symptomExists bool = false
	for _, ev := range events {
		for _, con := range ev.Conditions {
			if con.Symptom.Name == "symptom" {
				symptomExists = true
			}
		}
	}
	assert.True(t, symptomExists)
	assert.GreaterOrEqual(t, len(cats), 1)
}
