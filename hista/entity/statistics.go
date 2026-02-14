package entity

type FoodResult struct {
	SymptomID uint  `gorm:"column:symptom_id"`
	Severity  int   `gorm:"column:severity"`
	Hours72   int   `gorm:"column:hours_72"`
	Hours24   int   `gorm:"column:hours_24"`
	Hours1    int   `gorm:"column:hours_1"`
	Count     int64 `gorm:"column:count"` // TODO: Remove! Should not be needed anymore
}

type FoodResults []FoodResult

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

func (res FoodResults) ToResponse(count int64) SymptomStatisticsResponse {
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

type SymptomsResult struct {
	IngredientID  uint
	FoodCondition string
	Hours72       int
	Hours24       int
	Hours1        int
	Count         int64
}

type SymptomResults []SymptomsResult

type StatisticsByFood struct {
	IngredientID  uint   `json:"ingredientId"`
	FoodCondition string `json:"foodCondition"`
	Hours72       int    `json:"hours72"`
	Hours24       int    `json:"hours24"`
	Hours1        int    `json:"hours1"`
	Count         int64  `json:"count"`
}

type FoodStatisticsResponse struct {
	Statistics []StatisticsByFood `json:"statistics"`
}

func (res SymptomResults) ToResponse() FoodStatisticsResponse {
	var stats = []StatisticsByFood{}
	for _, res := range res {
		stat := StatisticsByFood(res)
		stats = append(stats, stat)
	}
	return FoodStatisticsResponse{Statistics: stats}
}

type CountResult struct {
	ID    uint  `gorm:"column:id"`
	Count int64 `gorm:"column:count"`
}
