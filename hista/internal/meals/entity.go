package meals

import (
	"time"

	"encore.dev/types/uuid"
)

type Freshness uint8

const (
	Fresh   Freshness = 0
	SameDay Freshness = 1
	Older   Freshness = 2
)

type Meal struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	PIID        uuid.UUID `gorm:"index;type:uuid;not null"`
	Date        time.Time
	Freshness   Freshness
	StressLevel uint8
	IsAlone     bool
	Foods       []Food `gorm:"constraint:OnDelete:CASCADE"`
}

func (m *Meal) SetPiid(id uuid.UUID) {
	m.PIID = id
}

type Meals []*Meal

type Food struct {
	ID           uint `gorm:"primaryKey"`
	Condition    FoodCondition
	Ingredient   Ingredient
	IngredientID uint
	MealID       uint
	Amount       *int
}

type Foods []Food

type FoodCondition string

const (
	Raw    FoodCondition = "raw"
	Cooked FoodCondition = "cooked"
)

type Ingredient struct {
	ID         uint      `gorm:"primaryKey"`
	PIID       uuid.UUID `gorm:"index;type:uuid;not null;uniqueIndex:idx_name_piid"`
	Name       string    `gorm:"uniqueIndex:idx_name_piid"`
	IsArchived bool      `gorm:"not null"`
	Nutrition  `gorm:"embeddedPrefix:nutrition_"`
	Items      []TemplateItem `gorm:"constraint:OnDelete:CASCADE"`
}

type Nutrition struct {
	Protein      *float32
	Carbohydrate *float32
	Fat          *float32
	Fiber        *float32
}

func (i *Ingredient) SetPiid(id uuid.UUID) {
	i.PIID = id
}

type Ingredients []*Ingredient

type Template struct {
	ID    uint           `gorm:"primaryKey"`
	PIID  uuid.UUID      `gorm:"type:uuid;index;not null;uniqueIndex:idx_template_name_piid"`
	Name  string         `gorm:"uniqueIndex:idx_template_name_piid"`
	Items []TemplateItem `gorm:"constraint:OnDelete:CASCADE"`
}

func (mt *Template) SetPiid(id uuid.UUID) {
	mt.PIID = id
}

type Templates []Template

type TemplateItem struct {
	ID           uint `gorm:"primaryKey"`
	Condition    FoodCondition
	TemplateID   uint
	IngredientID uint
}

type NutritionStatistic struct {
	Date      time.Time
	Nutrition Nutrition `gorm:"embedded"`
}

type NutritionStatistics []NutritionStatistic
