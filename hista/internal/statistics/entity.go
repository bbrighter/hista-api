package statistics

type FoodResult struct {
	SymptomID uint `gorm:"column:symptom_id"`
	Severity  int  `gorm:"column:severity"`
	Hours72   int  `gorm:"column:hours_72"`
	Hours24   int  `gorm:"column:hours_24"`
	Hours1    int  `gorm:"column:hours_1"`
}

type FoodResults []FoodResult
