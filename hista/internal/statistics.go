package internal

import (
	"context"
	"fmt"
	"time"

	"encore.app/hista/entity"
)

type (
	IStatisticsRepo interface {
		SymptomsAfterIngredients(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientId uint) (entity.FoodResults, error)
		CountMealsWithIngredients(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientId uint) (int64, error)
	}
	INutritionRepo interface {
		SelectAggregatedNutrition(ctx context.Context, truncateUnit string, from *time.Time, to *time.Time) (entity.NutritionStatistics, error)
	}

	IStatisticsUseCase interface {
		FindSymptomsForFoods(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientId uint) (entity.FoodResults, int64, error)
		FindNutrition(ctx context.Context, truncateUnit string, from *time.Time, to *time.Time) (entity.NutritionStatistics, error)
	}
)

type StatisticsUseCase struct {
	s IStatisticsRepo
	n INutritionRepo
}

func NewStatisticsUseCase(s IStatisticsRepo, n INutritionRepo) StatisticsUseCase {
	return StatisticsUseCase{s: s, n: n}
}

func (uc StatisticsUseCase) FindSymptomsForFoods(ctx context.Context, fromDate, toDate time.Time, ingredientId uint) (entity.FoodResults, int64, error) {
	counts, err := uc.s.CountMealsWithIngredients(ctx, fromDate, toDate, ingredientId)
	fmt.Printf("%v", counts)
	if err != nil {
		return entity.FoodResults{}, 0, err
	}
	results, err := uc.s.SymptomsAfterIngredients(ctx, fromDate, toDate, ingredientId)
	if err != nil {
		return entity.FoodResults{}, 0, err
	}

	// TODO: Nicht richtig. counts zählt pro Ingredient, results pro SymptomId. Nochmal nachdenken, was ich will!

	return results, counts, nil
}

func (uc StatisticsUseCase) FindNutrition(ctx context.Context, truncateUnit string, from *time.Time, to *time.Time) (entity.NutritionStatistics, error) {
	return uc.n.SelectAggregatedNutrition(ctx, truncateUnit, from, to)
}
