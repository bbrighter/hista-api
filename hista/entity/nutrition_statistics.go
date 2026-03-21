package entity

import (
	"slices"
	"time"
)

type NutritionStatistic struct {
	Date      time.Time
	Nutrition Nutrition `gorm:"embedded"`
}

type NutritionStatistics []NutritionStatistic

type NutritionStatisticResponse struct {
	Date      time.Time      `json:"time"`
	Nutrition *NutritionResp `json:"nutrition"`
}

type NutritionStatisticsResponse struct {
	Statistics []NutritionStatisticResponse `json:"statistics"`
}

func (ns NutritionStatistics) ToResp() NutritionStatisticsResponse {
	var statistics = []NutritionStatisticResponse{}
	for _, n := range ns {
		statistics = append(statistics, NutritionStatisticResponse{
			Date:      n.Date,
			Nutrition: n.Nutrition.toNutritionResp(),
		})
	}
	slices.SortFunc(statistics, func(a, b NutritionStatisticResponse) int {
		return a.Date.Compare(b.Date)
	})
	return NutritionStatisticsResponse{Statistics: statistics}
}
