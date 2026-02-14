package hista

import (
	"context"
	"time"

	"encore.app/errors"
	"encore.app/hista/entity"
	uuid "encore.dev/types/uuid"
)

type MealStatisticsParams struct {
	ID       uint      `json:"id"`
	FromDate time.Time `json:"fromDate"`
	ToDate   time.Time `json:"toDate"`
}

// encore:api auth method=GET path=/piid/:piid/statistics/ingredients
func (service *Service) GetStatisticsByIngredientId(ctx context.Context, piid uuid.UUID, params MealStatisticsParams) (resp entity.SymptomStatisticsResponse, err error) {
	if params.FromDate.IsZero() {
		return resp, errors.ErrorAttributeMustBeSet("fromDate")
	}
	if params.ToDate.IsZero() {
		return resp, errors.ErrorAttributeMustBeSet("toDate")
	}
	if params.ID == 0 {
		return resp, errors.ErrorAttributeMustBeSet("ids")
	}
	stats, count, err := service.statistics.FindSymptomsForFoods(ctx, params.FromDate, params.ToDate, params.ID)
	return stats.ToResponse(count), err
}
