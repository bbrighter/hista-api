package internal

import (
	"time"

	"encore.app/entity"
)

type StatisticsUseCase struct {
	repo IStatisticsRepo
}

func NewStatisticsUseCase(repo IStatisticsRepo) StatisticsUseCase {
	return StatisticsUseCase{repo: repo}
}

func (uc StatisticsUseCase) FindSymptomsForFoods(fromDate time.Time, toDate time.Time, ingredientIds []uint) (results entity.FoodResults, err error) {
	results, err = uc.repo.FindSymptomsForFoods(fromDate, toDate, ingredientIds)
	if err != nil {
		return results, err
	}
	var ids []uint
	for _, r := range results {
		ids = append(ids, r.SymptomID)
	}
	counts := uc.repo.CountSymptoms(ids)
	for i, r := range results {
		for _, count := range counts {
			if count.ID == r.SymptomID {
				results[i].Count = count.Count
			}
		}
	}
	return results, err
}
func (uc StatisticsUseCase) FindFoodForSymptoms(fromDate time.Time, toDate time.Time, symptomIds []uint) (results entity.SymptomResults, err error) {
	results, err = uc.repo.FindFoodForSymptoms(fromDate, toDate, symptomIds)
	if err != nil {
		return results, err
	}
	var ids []uint
	for _, r := range results {
		ids = append(ids, r.IngredientID)
	}
	counts := uc.repo.CountFoods(ids)
	for i, r := range results {
		for _, count := range counts {
			if count.ID == r.IngredientID {
				results[i].Count = count.Count
			}
		}
	}
	return results, err
}
