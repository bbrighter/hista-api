package status

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToSnakeCase(t *testing.T) {
	assert.Equal(t, "attribute", toSnakeCase("Attribute"))
	assert.Equal(t, "simple_attribute", toSnakeCase("SimpleAttribute"))
}

func TestUpdateStatusToMap(t *testing.T) {
	two := 2
	three := 3
	date := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	crash := false
	tests := map[string]struct {
		input  UpdateStatusParams
		result map[string]any
	}{
		"only one filled": {input: UpdateStatusParams{MorningFitness: &three}, result: map[string]any{"morning_fitness": 3}},
		"all filled": {input: UpdateStatusParams{
			MorningFitness:        &three,
			MorningSleep:          &three,
			DayFitness:            &three,
			EveningFitness:        &three,
			Depressive:            &two,
			Tense:                 &two,
			MoodSwings:            &two,
			Irritable:             &two,
			LossOfInterest:        &two,
			ConcentrationProblems: &two,
			LackOfDrive:           &two,
			AppetiteChanges:       &two,
			SleepProblems:         &two,
			Overwhelmed:           &two,
			Date:                  &date,
			Crash:                 &crash,
		}, result: map[string]any{
			"morning_fitness":        3,
			"morning_sleep":          3,
			"day_fitness":            3,
			"evening_fitness":        3,
			"depressive":             2,
			"tense":                  2,
			"mood_swings":            2,
			"irritable":              2,
			"loss_of_interest":       2,
			"concentration_problems": 2,
			"lack_of_drive":          2,
			"appetite_changes":       2,
			"sleep_problems":         2,
			"overwhelmed":            2,
			"date":                   time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
			"crash":                  false,
		}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.result, test.input.ToMap())
		})
	}
}
