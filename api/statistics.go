package api

import (
	"context"
	"log"
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
	if len(params.IDs) == 0 {
		return resp, errors.ErrorAttributeMustBeSet("ids")
	}
	stats, err := service.statistics.FindFoodForSymptoms(ctx, params.FromDate, params.ToDate, params.IDs)
	if len(stats) > 0 {
		log.Printf("Stats %v", stats[0].Count)
	}
	return stats.ToResponse(), err
}

// encore:api auth method=GET path=/statistics/ingredients
func (service *Service) GetStatisticsByIngredientsIds(ctx context.Context, params StatisticParams) (resp entity.SymptomStatisticsResponse, err error) {
	if params.FromDate.IsZero() {
		return resp, errors.ErrorAttributeMustBeSet("fromDate")
	}
	if params.ToDate.IsZero() {
		return resp, errors.ErrorAttributeMustBeSet("toDate")
	}
	if len(params.IDs) == 0 {
		return resp, errors.ErrorAttributeMustBeSet("ids")
	}
	stats, err := service.statistics.FindSymptomsForFoods(ctx, params.FromDate, params.ToDate, params.IDs)
	return stats.ToResponse(), err
}
