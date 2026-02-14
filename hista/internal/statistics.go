package internal

import (
	"context"
	"fmt"
	"time"

	"encore.app/errors"
	"encore.app/hista/entity"
)

type (
	IStatisticsRepo interface {
		SymptomsAfterIngredients(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientId uint) (entity.FoodResults, error)
		CountMealsWithIngredients(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientId uint) (int64, error)
	}
	IOldStatisticsRepo interface {
		FindSymptomsForFoods(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientIds []uint) (entity.FoodResults, error)
		FindFoodForSymptoms(ctx context.Context, fromDate time.Time, toDate time.Time, symptomIds []uint) (entity.SymptomResults, error)
		CountFoods(ctx context.Context, symptomIds []uint) ([]entity.CountResult, error)
		CountSymptoms(ctx context.Context, ingredientIds []uint) ([]entity.CountResult, error)
	}

	IStatisticsUseCase interface {
		FindSymptomsForFoods(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientId uint) (entity.FoodResults, int64, error)
		FindFoodForSymptoms(ctx context.Context, fromDate time.Time, toDate time.Time, symptomIds []uint) (entity.SymptomResults, error)
	}
)

type StatiaticsUseCase struct {
	r IStatisticsRepo
}

func NewStatisticsUseCase(repo IStatisticsRepo) StatiaticsUseCase {
	return StatiaticsUseCase{r: repo}
}

func (uc StatiaticsUseCase) FindSymptomsForFoods(ctx context.Context, fromDate, toDate time.Time, ingredientId uint) (entity.FoodResults, int64, error) {
	counts, err := uc.r.CountMealsWithIngredients(ctx, fromDate, toDate, ingredientId)
	fmt.Printf("%v", counts)
	if err != nil {
		return entity.FoodResults{}, 0, err
	}
	results, err := uc.r.SymptomsAfterIngredients(ctx, fromDate, toDate, ingredientId)
	if err != nil {
		return entity.FoodResults{}, 0, err
	}

	// TODO: Nicht richtig. counts zählt pro Ingredient, results pro SymptomId. Nochmal nachdenken, was ich will!

	return results, counts, nil
}

func (uc StatiaticsUseCase) FindFoodForSymptoms(ctx context.Context, fromDate time.Time, toDate time.Time, symptomIds []uint) (entity.SymptomResults, error) {
	return entity.SymptomResults{}, nil
}

type StatisticsOldUseCase struct {
	repo IOldStatisticsRepo
}

func NewStatisticsOldUseCase(repo IOldStatisticsRepo) StatisticsOldUseCase {
	return StatisticsOldUseCase{repo: repo}
}

func (uc StatisticsOldUseCase) FindSymptomsForFoods(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientIds []uint) (results entity.FoodResults, err error) {
	results, err = uc.repo.FindSymptomsForFoods(ctx, fromDate, toDate, ingredientIds)
	if err != nil {
		return results, errors.MapError(err)
	}
	var ids []uint
	for _, r := range results {
		ids = append(ids, r.SymptomID)
	}
	counts, err := uc.repo.CountSymptoms(ctx, ids)
	if err != nil {
		return entity.FoodResults{}, errors.MapError(err)
	}
	for i, r := range results {
		for _, count := range counts {
			if count.ID == r.SymptomID {
				results[i].Count = count.Count
			}
		}
	}
	return results, errors.MapError(err)
}
func (uc StatisticsOldUseCase) FindFoodForSymptoms(ctx context.Context, fromDate time.Time, toDate time.Time, symptomIds []uint) (results entity.SymptomResults, err error) {
	results, err = uc.repo.FindFoodForSymptoms(ctx, fromDate, toDate, symptomIds)
	if err != nil {
		return results, errors.MapError(err)
	}
	var ids []uint
	for _, r := range results {
		ids = append(ids, r.IngredientID)
	}
	counts, err := uc.repo.CountFoods(ctx, ids)
	if err != nil {
		return entity.SymptomResults{}, errors.MapError(err)
	}
	for i, r := range results {
		for _, count := range counts {
			if count.ID == r.IngredientID {
				results[i].Count = count.Count
			}
		}
	}
	return results, errors.MapError(err)
}
