package statistics

import (
	"context"
	"time"

	"encore.app/hista/internal/meals"
	"gorm.io/gorm"
)

type StatisticsService struct {
	m *meals.MealRepository
	s *statisticsRepo
}

func NewStatisticsService(db *gorm.DB) *StatisticsService {
	m := meals.NewMealRepository(db)
	s := newStatisticsRepo(db)
	return &StatisticsService{m: m, s: s}
}

func (s StatisticsService) AggregateNutrition(
	ctx context.Context,
	truncateUnit string,
	from, to *time.Time) (meals.NutritionStatistics, error) {
	return s.m.AggregateNutrition(ctx, truncateUnit, from, to)
}
func (s StatisticsService) SymptomsAfterIngredient(ctx context.Context, from time.Time, to time.Time, ingredientId uint) (FoodResults, int64, error) {
	counts, err := s.s.CountMealsWithIngredients(ctx, from, to, ingredientId)
	if err != nil {
		return FoodResults{}, 0, err
	}
	results, err := s.s.SymptomsAfterIngredients(ctx, from, to, ingredientId)
	if err != nil {
		return FoodResults{}, 0, err
	}
	return results, counts, nil
}
func (s StatisticsService) DiaryEntries(ctx context.Context) {
	print("ok")
}
