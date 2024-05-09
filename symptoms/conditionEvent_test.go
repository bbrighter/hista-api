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
