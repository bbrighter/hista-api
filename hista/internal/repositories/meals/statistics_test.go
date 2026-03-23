package meals

import (
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
)

func (s *MealRepoTestSuite) TestSelectAggregatedNutrition() {
	var (
		proteinVal float32 = 10.0
		fiberVal   float32 = 2.0
		carbsVal   float32 = 5.2
		fatVal     float32 = 1.0
		amountVal  int     = 95
	)

	meal1Date := time.Date(2022, 3, 3, 12, 0, 0, 0, time.UTC)
	meal1 := entity.Meal{Date: meal1Date}
	meal2 := entity.Meal{Date: meal1Date.Add(23 * time.Hour)}
	generic_queries.Create(s.ctx, s.tx, &meal1)
	generic_queries.Create(s.ctx, s.tx, &meal2)

	ingredientWithNutrition := entity.Ingredient{
		Name: "name",
		Nutrition: entity.Nutrition{
			Protein:      &proteinVal,
			Fiber:        &fiberVal,
			Fat:          &fatVal,
			Carbohydrate: &carbsVal,
		},
	}
	ingredientWithoutNutrition := entity.Ingredient{Name: "no nutrition"}
	generic_queries.Create(s.ctx, s.tx, &ingredientWithNutrition)
	generic_queries.Create(s.ctx, s.tx, &ingredientWithoutNutrition)

	foodWithAmount := entity.Food{IngredientID: ingredientWithNutrition.ID, MealID: meal1.ID, Amount: &amountVal}
	food2 := entity.Food{IngredientID: ingredientWithoutNutrition.ID, MealID: meal1.ID}
	food3 := entity.Food{IngredientID: ingredientWithNutrition.ID, MealID: meal2.ID}
	generic_queries.Create(s.ctx, s.tx, &foodWithAmount)
	generic_queries.Create(s.ctx, s.tx, &food2)
	generic_queries.Create(s.ctx, s.tx, &food3)

	tests := map[string]struct {
		fromDate      *time.Time
		toDate        *time.Time
		interval      string
		expectedError bool
		expectedLen   int
	}{
		"no dates":           {interval: "hour", expectedLen: 1},
		"invalid interval":   {interval: "something", expectedError: true},
		"from date":          {interval: "hour", fromDate: new(meal1Date.Add(-time.Minute)), expectedLen: 1},
		"from date too late": {interval: "hour", fromDate: new(meal1Date.Add(time.Minute)), expectedLen: 0},
		"to date":            {interval: "hour", toDate: new(meal1Date.Add(24 * time.Hour)), expectedLen: 1},
		"to date too early":  {interval: "hour", toDate: new(meal1Date.Add(-24 * time.Hour)), expectedLen: 0},
	}
	for name, test := range tests {
		s.Run(name, func() {
			nutrition, err := s.repo.SelectAggregatedNutrition(s.ctx, test.interval, test.fromDate, test.toDate)
			if test.expectedError {
				s.Error(err)
				return
			}
			s.Len(nutrition, test.expectedLen)
			if test.expectedLen == 0 {
				return
			}

			scale := float32(amountVal) / 100.0
			expectedProtein := proteinVal * scale
			expectedFat := fatVal * scale
			expectedFiber := fiberVal * scale
			expectedCarbs := carbsVal * scale

			s.Equal(meal1Date, nutrition[0].Date.UTC())
			s.InDelta(expectedProtein, *nutrition[0].Nutrition.Protein, 0.001)
			s.InDelta(expectedFat, *nutrition[0].Nutrition.Fat, 0.001)
			s.InDelta(expectedFiber, *nutrition[0].Nutrition.Fiber, 0.001)
			s.InDelta(expectedCarbs, *nutrition[0].Nutrition.Carbohydrate, 0.001)
		})
	}
}

func (s *MealRepoTestSuite) TestSelectAggregatedNutritionInvalidInput() {
	_, err := s.repo.SelectAggregatedNutrition(s.ctx, "one and a half hours", nil, nil)
	s.Error(err)
	s.ErrorContains(err, "invalid truncateUnit")
}
