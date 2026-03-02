package entity

import (
	"sort"
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
	PIID        uuid.UUID `gorm:"type:uuid;primaryKey"`
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

type MealMetaResponse struct {
	ID   uint      `json:"id"`
	Date time.Time `json:"date"`
}

type MealResponse struct {
	ID          uint           `json:"id"`
	Date        time.Time      `json:"date"`
	Freshness   Freshness      `json:"freshness"`
	StressLevel uint8          `json:"stressLevel"`
	IsAlone     bool           `json:"isAlone"`
	Foods       []FoodResponse `json:"foods"`
}

type MealsResponse struct {
	Meals []MealMetaResponse `json:"meals"`
}

type IDResponse struct {
	ID uint `json:"id"`
}

func (meal Meal) ToMealMetaResponse() MealMetaResponse {
	return MealMetaResponse{
		ID:   meal.ID,
		Date: meal.Date,
	}
}

func (meal Meal) ToMealResponse() MealResponse {
	var resp = MealResponse{
		ID:          meal.ID,
		Date:        meal.Date,
		Foods:       NonPtFoods(meal.Foods).ToFoodsResponse(),
		Freshness:   meal.Freshness,
		StressLevel: meal.StressLevel,
		IsAlone:     meal.IsAlone,
	}
	return resp
}

func (meals Meals) ToMealsResponse() MealsResponse {
	var resps []MealMetaResponse
	for _, m := range meals {
		resps = append(resps, m.ToMealMetaResponse())
	}
	sort.Slice(resps, func(i, j int) bool {
		return resps[i].Date.Sub(resps[j].Date) > 0
	})
	return MealsResponse{resps}
}

type PostMealParams struct {
	Date time.Time `json:"date"`
}

type PatchMealParams struct {
	Date        *time.Time `json:"date" encore:"optional"`
	Freshness   *Freshness `json:"freshness" encore:"optional"`
	StressLevel *uint8     `json:"stressLevel" encore:"optional"`
	IsAlone     *bool      `json:"isAlone" encore:"optional"`
}
