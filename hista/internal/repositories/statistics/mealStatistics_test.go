package statistics

import (
	"time"

	"encore.app/hista/entity"
)

func (s *StatsRepoTestSuite) TestSymptomsAfterIngredients() {
	validFromDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := map[string]struct {
		fromDate      time.Time
		toDate        time.Time
		ingredientIds uint
		want          []entity.FoodResult
	}{
		"ok": {
			fromDate:      validFromDate,
			toDate:        validFromDate.Add(time.Hour * 24 * 365),
			ingredientIds: 1,
			want: []entity.FoodResult{
				{SymptomID: 1, Hours1: 0, Hours24: 1, Hours72: 1, Severity: 2},
				{SymptomID: 2, Hours1: 0, Hours24: 0, Hours72: 1, Severity: 2},
			},
		},
		"outside date range": {
			fromDate:      validFromDate.Add(time.Hour * 24 * 365),
			toDate:        validFromDate.Add(time.Hour * 24 * 366),
			ingredientIds: 1,
			want:          []entity.FoodResult{},
		},
		"no ingredient": {
			fromDate:      validFromDate,
			toDate:        validFromDate.Add(time.Hour * 24 * 365),
			ingredientIds: 100,
			want:          []entity.FoodResult{},
		},
	}
	for name, tt := range tests {
		s.Run(name, func() {
			got, gotErr := s.stats.SymptomsAfterIngredients(s.ctx, tt.fromDate, tt.toDate, tt.ingredientIds)
			s.NoError(gotErr)
			s.ElementsMatch(got, tt.want)

		})
	}
}

func (s *StatsRepoTestSuite) TestCountMealsWithIngredients() {
	validFromDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := map[string]struct {
		fromDate      time.Time
		toDate        time.Time
		ingredientIds uint
		want          int64
	}{
		"ok": {
			fromDate:      validFromDate,
			toDate:        validFromDate.Add(time.Hour * 24 * 365),
			ingredientIds: 1,
			want:          2,
		},
		"outside date range": {
			fromDate:      validFromDate.Add(time.Hour * 24 * 365),
			toDate:        validFromDate.Add(time.Hour * 24 * 366),
			ingredientIds: 1,
			want:          0,
		},
		"no ingredient": {
			fromDate:      validFromDate,
			toDate:        validFromDate.Add(time.Hour * 24 * 365),
			ingredientIds: 100,
			want:          0,
		},
	}
	for name, tt := range tests {
		s.Run(name, func() {
			got, gotErr := s.stats.CountMealsWithIngredients(s.ctx, tt.fromDate, tt.toDate, tt.ingredientIds)
			s.NoError(gotErr)
			s.EqualValues(got, tt.want)

		})
	}
}
