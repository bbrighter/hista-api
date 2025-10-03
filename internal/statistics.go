package internal

import (
	"context"
	"time"

	"encore.app/entity"
)

type (
	IStatisticsRepo interface {
		FindSymptomsForFoods(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientIds []uint) (entity.FoodResults, error)
		FindFoodForSymptoms(ctx context.Context, fromDate time.Time, toDate time.Time, symptomIds []uint) (entity.SymptomResults, error)
		CountFoods(ctx context.Context, symptomIds []uint) ([]entity.CountResult, error)
		CountSymptoms(ctx context.Context, ingredientIds []uint) ([]entity.CountResult, error)
	}

	IStatisticsUseCase interface {
		FindSymptomsForFoods(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientIds []uint) (entity.FoodResults, error)
		FindFoodForSymptoms(ctx context.Context, fromDate time.Time, toDate time.Time, symptomIds []uint) (entity.SymptomResults, error)
	}
)

type StatisticsUseCase struct {
	repo IStatisticsRepo
}

func NewStatisticsUseCase(repo IStatisticsRepo) StatisticsUseCase {
	return StatisticsUseCase{repo: repo}
}

func (uc StatisticsUseCase) FindSymptomsForFoods(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientIds []uint) (results entity.FoodResults, err error) {
	results, err = uc.repo.FindSymptomsForFoods(ctx, fromDate, toDate, ingredientIds)
	if err != nil {
		return results, errorMapper(err)
	}
	var ids []uint
	for _, r := range results {
		ids = append(ids, r.SymptomID)
	}
	counts, err := uc.repo.CountSymptoms(ctx, ids)
	if err != nil {
		return entity.FoodResults{}, errorMapper(err)
	}
	for i, r := range results {
		for _, count := range counts {
			if count.ID == r.SymptomID {
				results[i].Count = count.Count
			}
		}
	}
	return results, errorMapper(err)
}
func (uc StatisticsUseCase) FindFoodForSymptoms(ctx context.Context, fromDate time.Time, toDate time.Time, symptomIds []uint) (results entity.SymptomResults, err error) {
	results, err = uc.repo.FindFoodForSymptoms(ctx, fromDate, toDate, symptomIds)
	if err != nil {
		return results, errorMapper(err)
	}
	var ids []uint
	for _, r := range results {
		ids = append(ids, r.IngredientID)
	}
	counts, err := uc.repo.CountFoods(ctx, ids)
	if err != nil {
		return entity.SymptomResults{}, errorMapper(err)
	}
	for i, r := range results {
		for _, count := range counts {
			if count.ID == r.IngredientID {
				results[i].Count = count.Count
			}
		}
	}
	return results, errorMapper(err)
}
