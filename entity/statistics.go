package entity

type FoodResult struct {
	SymptomID uint
	Severity  int
	Hours72   int
	Hours24   int
	Hours1    int
	Count     int64
}

type FoodResults []FoodResult

type StatisticBySymptom struct {
	SymptomID uint  `json:"symptomId"`
	Severity  int   `json:"severity"`
	Hours72   int   `json:"hours72"`
	Hours24   int   `json:"hours24"`
	Hours1    int   `json:"hours1"`
	Count     int64 `json:"count"`
}

type SymptomStatisticsResponse struct {
	Statistics []StatisticBySymptom `json:"statistics"`
}

func (res FoodResults) ToResponse() SymptomStatisticsResponse {
	var stats = []StatisticBySymptom{}
	for _, res := range res {
		var stat = StatisticBySymptom{
			SymptomID: res.SymptomID,
			Severity:  res.Severity,
			Hours72:   res.Hours72,
			Hours24:   res.Hours24,
			Hours1:    res.Hours1,
		}
		stats = append(stats, stat)
	}
	return SymptomStatisticsResponse{Statistics: stats}
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
		var stat = StatisticsByFood{
			IngredientID:  res.IngredientID,
			FoodCondition: res.FoodCondition,
			Hours72:       res.Hours72,
			Hours24:       res.Hours24,
			Hours1:        res.Hours1,
		}
		stats = append(stats, stat)
	}
	return FoodStatisticsResponse{Statistics: stats}
}

type CountResult struct {
	ID    uint
	Count int64
}
