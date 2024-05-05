package states

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewState(t *testing.T) {
	var state *State = newState(time.Now())

	assert.GreaterOrEqual(t, time.Now(), state.Date)
}

func TestCreateState(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var state *State
	var err error

	err = state.create(service)
	assert.Error(t, err)

	state = newState(time.Now())
	err = state.create(service)
	assert.NoError(t, err)
}

func TestGetStates(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var states States
	states = getStates(service)
	assert.Len(t, states, 0)

	var state State = service.testCreateState(t)
	states = getStates(service)
	assert.Len(t, states, 1)
	assert.Equal(t, state.ID, states[0].ID)
	assert.False(t, states[0].Date.IsZero(), "date is set; no comparison because time.Now is used")
}

func TestGetState(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	service.testCreateState(t)
	var state State
	var err error
	state, err = getState(service, 1)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, state.ID)

	// Not found
	_, err = getState(service, 1000)
	assert.Error(t, err)
}
