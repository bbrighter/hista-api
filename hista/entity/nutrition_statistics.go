package entity

import "time"

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
	return NutritionStatisticsResponse{Statistics: statistics}
}
