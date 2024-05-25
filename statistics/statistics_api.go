package statistics

import (
	"context"
	"time"

	"encore.app/errors"
)

type StatisticParams struct {
	SymptomIDs []uint    `json:"symptomIds"`
	FromDate   time.Time `json:"fromDate"`
	ToDate     time.Time `json:"toDate"`
}

// encore:api auth method=GET path=/statistics
func (service *Service) GetSymptomsBySymptomIDs(ctx context.Context, params StatisticParams) (StatisticsResponse, error) {
	var resp = StatisticsResponse{}
	print("These are the params:", params.FromDate.String(), params.ToDate.String(), params.SymptomIDs)
	if params.FromDate.IsZero() {
		return resp, errors.ErrorAttributeMustBeSet("fromDate")
	}
	if params.ToDate.IsZero() {
		return resp, errors.ErrorAttributeMustBeSet("toDate")
	}
	if params.SymptomIDs == nil || len(params.SymptomIDs) == 0 {
		return resp, errors.ErrorAttributeMustBeSet("symptomIds")
	}
	resp, err := findFoodForSymptoms(service, params.FromDate, params.ToDate, params.SymptomIDs)
	if err != nil {
		return resp, err
	}
	return resp, err
}
