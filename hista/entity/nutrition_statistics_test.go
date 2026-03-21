package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNutritionStatisticsToResp(t *testing.T) {
	var protein float32 = 100
	tests := []struct {
		name     string
		input    NutritionStatistics
		expected []time.Time
	}{
		{
			name:     "empty statistics",
			input:    NutritionStatistics{},
			expected: []time.Time{},
		},
		{
			name: "single statistic",
			input: NutritionStatistics{
				{
					Date:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Nutrition: Nutrition{Protein: &protein},
				},
			},
			expected: []time.Time{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			name: "multiple statistics sorted descending by date",
			input: NutritionStatistics{
				{
					Date:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Nutrition: Nutrition{Protein: &protein},
				},
				{
					Date:      time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
					Nutrition: Nutrition{Protein: &protein},
				},
				{
					Date:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					Nutrition: Nutrition{Protein: &protein},
				},
			},
			expected: []time.Time{
				time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToResp()

			assert.Len(t, result.Statistics, len(tt.expected))

			for i, expectedDate := range tt.expected {
				assert.Equal(t, expectedDate, result.Statistics[i].Date)
			}
		})
	}
}
