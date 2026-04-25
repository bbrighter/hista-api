package hista

import (
	"context"
	"slices"
	"time"

	"encore.app/errors"
	"encore.app/hista/internal/meals"
	"encore.app/hista/internal/statistics"
	"encore.dev/types/option"
	uuid "encore.dev/types/uuid"
)

type MealStatisticsParams struct {
	ID       uint      `json:"id"`
	FromDate time.Time `json:"fromDate"`
	ToDate   time.Time `json:"toDate"`
}

type StatisticBySymptom struct {
	SymptomID uint `json:"symptomId"`
	Severity  int  `json:"severity"`
	Hours72   int  `json:"hours72"`
	Hours24   int  `json:"hours24"`
	Hours1    int  `json:"hours1"`
}

type SymptomStatisticsResponse struct {
	Count      int64                `json:"count"`
	Statistics []StatisticBySymptom `json:"statistics"`
}

func toSymptomStatisticsResponse(res statistics.FoodResults, count int64) SymptomStatisticsResponse {
	var stats = []StatisticBySymptom{}
	for _, res := range res {
		stats = append(stats, StatisticBySymptom{
			SymptomID: res.SymptomID,
			Severity:  res.Severity,
			Hours72:   res.Hours72,
			Hours24:   res.Hours24,
			Hours1:    res.Hours1,
		})
	}
	return SymptomStatisticsResponse{Statistics: stats, Count: count}
}

// encore:api auth method=GET path=/piid/:piid/statistics/ingredients
func (service *Service) GetStatisticsByIngredientId(ctx context.Context, piid uuid.UUID, params MealStatisticsParams) (resp SymptomStatisticsResponse, err error) {
	if params.FromDate.IsZero() {
		return resp, errors.ErrorAttributeMustBeSet("fromDate")
	}
	if params.ToDate.IsZero() {
		return resp, errors.ErrorAttributeMustBeSet("toDate")
	}
	if params.ID == 0 {
		return resp, errors.ErrorAttributeMustBeSet("ids")
	}
	stats, count, err := service.stats.SymptomsAfterIngredient(ctx, params.FromDate, params.ToDate, params.ID)
	if err != nil {
		return resp, errors.MapError(err)
	}
	return toSymptomStatisticsResponse(stats, count), nil
}

type NutritionResp struct {
	Protein      float32 `json:"protein"`
	Carbohydrate float32 `json:"carbohydrate"`
	Fat          float32 `json:"fat"`
	Fiber        float32 `json:"fiber"`
}

type NutritionStatisticResponse struct {
	Date      time.Time      `json:"time"`
	Nutrition *NutritionResp `json:"nutrition"`
}

type NutritionStatisticsResponse struct {
	Statistics []NutritionStatisticResponse `json:"statistics"`
}

func toNutritionResp(n meals.Nutrition) *NutritionResp {
	if n.Protein == nil || n.Fat == nil || n.Fiber == nil || n.Carbohydrate == nil {
		return nil
	}
	return &NutritionResp{
		Protein:      *n.Protein,
		Carbohydrate: *n.Carbohydrate,
		Fat:          *n.Fat,
		Fiber:        *n.Fiber,
	}
}

func toNutritionStatisticsResponse(ns meals.NutritionStatistics) NutritionStatisticsResponse {
	var statistics = []NutritionStatisticResponse{}
	for _, n := range ns {
		statistics = append(statistics, NutritionStatisticResponse{
			Date:      n.Date,
			Nutrition: toNutritionResp(n.Nutrition),
		})
	}
	slices.SortFunc(statistics, func(a, b NutritionStatisticResponse) int {
		return b.Date.Compare(a.Date)
	})
	return NutritionStatisticsResponse{Statistics: statistics}
}

type NutritionStatisticsParams struct {
	Interval string                   `query:"interval"`
	From     option.Option[time.Time] `query:"from"`
	To       option.Option[time.Time] `query:"to"`
}

// encore:api auth method=GET path=/piid/:piid/statistics/nutrition
func (s *Service) GetNutritionByInterval(ctx context.Context, piid uuid.UUID, params NutritionStatisticsParams) (NutritionStatisticsResponse, error) {
	from := params.From.PtrOrNil()
	to := params.To.PtrOrNil()
	nutrition, err := s.stats.AggregateNutrition(ctx, params.Interval, from, to)
	if err != nil {
		return NutritionStatisticsResponse{}, errors.MapError(err)
	}
	return toNutritionStatisticsResponse(nutrition), nil
}
