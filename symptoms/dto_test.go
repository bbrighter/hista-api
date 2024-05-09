package symptoms

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSymptomToResponse(t *testing.T) {
	t.Parallel()
	var symptom = Symptom{
		ID:                1,
		Name:              "Name",
		SymptomCategoryID: 3,
	}
	resp := symptom.toResponse()
	assert.EqualValues(t, 1, resp.ID)
	assert.Equal(t, "Name", resp.Name)
	assert.EqualValues(t, 3, resp.CategoryID)
}

func TestSymptomCategoriesToResponse(t *testing.T) {
	t.Parallel()
	var cats = SymptomCategories{
		SymptomCategory{
			ID:   1,
			Name: "Cat",
			Symptoms: []Symptom{
				{ID: 10, Name: "Sym", SymptomCategoryID: 1},
			},
		},
	}
	resp := cats.toResponse()
	assert.Len(t, resp.Categories, 1)
	respCat := resp.Categories[0]
	assert.EqualValues(t, 1, respCat.ID)
	assert.Equal(t, "Cat", respCat.Name)
	assert.Len(t, respCat.Symptoms, 1)
	respSym := respCat.Symptoms[0]
	assert.EqualValues(t, 10, respSym.ID)
	assert.EqualValues(t, 1, respSym.CategoryID)
	assert.Equal(t, "Sym", respSym.Name)
}

func TestConditionToResponse(t *testing.T) {
	t.Parallel()
	var condition = Condition{
		ID:               1,
		Severity:         High,
		ConditionEventID: 2,
		SymptomID:        3,
		Symptom: Symptom{
			ID:                3,
			Name:              "Sym",
			SymptomCategoryID: 10,
		},
	}
	resp := condition.toResponse()
	assert.EqualValues(t, 1, resp.ID)
	assert.Equal(t, High, resp.Severity)
	assert.EqualValues(t, 3, resp.Symptom.ID)
	assert.Equal(t, "Sym", resp.Symptom.Name)
	assert.EqualValues(t, 10, resp.Symptom.CategoryID)
}

func TestConditionEventsToResponse(t *testing.T) {
	t.Parallel()
	var events = ConditionEvents{
		ConditionEvent{
			ID:   1,
			Date: time.Now(),
			Conditions: []Condition{
				{ID: 1},
			},
		},
		ConditionEvent{
			ID:         2,
			Date:       time.Now().Add(time.Hour),
			Conditions: []Condition{},
		},
	}
	resp := events.toResponse()
	assert.Len(t, resp.ConditionEvents, 2)
	resp1 := resp.ConditionEvents[0]
	assert.EqualValues(t, 2, resp1.ID)
}

func TestConditionEventToResponse(t *testing.T) {
	t.Parallel()
	var event = ConditionEvent{
		ID:   1,
		Date: time.Now(),
		Conditions: []Condition{
			{
				ID:               10,
				SymptomID:        5,
				Severity:         High,
				ConditionEventID: 1,
				Symptom:          Symptom{ID: 5, Name: "Sym", SymptomCategoryID: 3},
			},
		},
	}
	resp := event.toResponse()
	assert.EqualValues(t, 1, resp.ID)
	assert.GreaterOrEqual(t, time.Now(), resp.Date)
	assert.Len(t, resp.Conditions, 1)
	cond := resp.Conditions[0]
	assert.EqualValues(t, 10, cond.ID)
	assert.EqualValues(t, 5, cond.Symptom.ID)
	assert.Equal(t, High, cond.Severity)
	assert.EqualValues(t, 3, cond.Symptom.CategoryID)
}
