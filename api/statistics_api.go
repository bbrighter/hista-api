package api

import (
	"context"
	"time"

	entity "encore.app/entity"
	"encore.app/errors"
)

type StatisticParams struct {
	IDs      []uint    `json:"ids"`
	FromDate time.Time `json:"fromDate"`
	ToDate   time.Time `json:"toDate"`
}

// encore:api auth method=GET path=/statistics/symptoms
func (service *Service) GetStatisticsBySymptomIds(ctx context.Context, params StatisticParams) (resp entity.FoodStatisticsResponse, err error) {
	if params.FromDate.IsZero() {
		return resp, errors.ErrorAttributeMustBeSet("fromDate")
	}
	if params.ToDate.IsZero() {
		return resp, errors.ErrorAttributeMustBeSet("toDate")
	}
	if params.IDs == nil || len(params.IDs) == 0 {
		return resp, errors.ErrorAttributeMustBeSet("ids")
	}
	stats, err := service.statistics.FindFoodForSymptoms(params.FromDate, params.ToDate, params.IDs)
	return stats.ToResponse(), err

	// resp, err := findFoodForSymptoms(service, params.FromDate, params.ToDate, params.IDs)
	// if err != nil {
	// 	return resp, err
	// }

	// var relevantIngredientIDs []uint
	// for _, result := range resp.Statistics {
	// 	relevantIngredientIDs = append(relevantIngredientIDs, result.IngredientID)
	// }
	// counts := countFoods(service, relevantIngredientIDs)
	// for i, result := range resp.Statistics {
	// 	for _, count := range counts {
	// 		if count.ID == result.IngredientID {
	// 			resp.Statistics[i].Count = count.Count
	// 		}
	// 	}
	// }

	// return resp, err
}

// encore:api auth method=GET path=/statistics/ingredients
func (service *Service) GetStatisticsByIngredientsIds(ctx context.Context, params StatisticParams) (resp entity.SymptomStatisticsResponse, err error) {
	if params.FromDate.IsZero() {
		return resp, errors.ErrorAttributeMustBeSet("fromDate")
	}
	if params.ToDate.IsZero() {
		return resp, errors.ErrorAttributeMustBeSet("toDate")
	}
	if params.IDs == nil || len(params.IDs) == 0 {
		return resp, errors.ErrorAttributeMustBeSet("ids")
	}
	stats, err := service.statistics.FindSymptomsForFoods(params.FromDate, params.ToDate, params.IDs)
	return stats.ToResponse(), err
	// resp, err := findSymptomsForFoods(service, params.FromDate, params.ToDate, params.IDs)
	// if err != nil {
	// 	return resp, err
	// }
	// var relevantSymptomIds []uint
	// for _, result := range resp.Statistics {
	// 	relevantSymptomIds = append(relevantSymptomIds, result.SymptomID)
	// }
	// counts := countSymptoms(service, relevantSymptomIds)
	// for i, result := range resp.Statistics {
	// 	for _, count := range counts {
	// 		if count.ID == result.SymptomID {
	// 			resp.Statistics[i].Count = count.Count
	// 		}
	// 	}
	// }
	// return resp, err
}
