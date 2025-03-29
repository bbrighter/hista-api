package entity

import (
	"sort"
	"time"
)

type ConditionEvent struct {
	ID         uint
	Date       time.Time
	Conditions []Condition
}

type Condition struct {
	ID               uint
	Symptom          Symptom
	SymptomID        uint
	Severity         Severity
	ConditionEventID uint
}

type Symptom struct {
	ID                uint
	Name              string
	SymptomCategoryID uint
}

type SymptomCategory struct {
	ID       uint
	Name     string
	Symptoms []Symptom
}

type ConditionEvents []ConditionEvent

type Conditions []Condition

type Severity uint8

const (
	VeryLowSeverity  Severity = 1
	LowSeverity      Severity = 2
	MediumSeverity   Severity = 3
	HighSeverity     Severity = 4
	VeryHighSeverity Severity = 5
)

type Symptoms []Symptom

type SymptomCategories []SymptomCategory

type ConditionEventsResponse struct {
	ConditionEvents []ConditionEventMetaResponse `json:"conditionEvents"`
}

type ConditionEventMetaResponse struct {
	ID   uint      `json:"id"`
	Date time.Time `json:"date"`
}

type SymptomResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	CategoryID uint   `json:"categoryId"`
}

type SymptomCategoryResponse struct {
	ID       uint              `json:"id"`
	Name     string            `json:"name"`
	Symptoms []SymptomResponse `json:"symptoms"`
}

type SymptomCategoriesResponse struct {
	Categories []SymptomCategoryResponse
}

type ConditionEventResponse struct {
	ID         uint                `json:"id"`
	Date       time.Time           `json:"date"`
	Conditions []ConditionResponse `json:"conditions"`
}

type ConditionResponse struct {
	ID       uint            `json:"id"`
	Symptom  SymptomResponse `json:"symptom"`
	Severity Severity        `json:"severity"`
}

func (symptom Symptom) ToResponse() SymptomResponse {
	return SymptomResponse{
		ID:         symptom.ID,
		Name:       symptom.Name,
		CategoryID: symptom.SymptomCategoryID,
	}
}

func (cats SymptomCategories) ToResponse() SymptomCategoriesResponse {
	var resp = []SymptomCategoryResponse{}
	for _, cat := range cats {
		var symptomsResp = []SymptomResponse{}
		for _, sym := range cat.Symptoms {
			symptomsResp = append(symptomsResp, sym.ToResponse())
		}
		var c = SymptomCategoryResponse{
			ID:       cat.ID,
			Name:     cat.Name,
			Symptoms: symptomsResp,
		}
		resp = append(resp, c)
	}
	return SymptomCategoriesResponse{Categories: resp}
}

func (condition Condition) ToResponse() ConditionResponse {
	return ConditionResponse{
		ID:       condition.ID,
		Symptom:  condition.Symptom.ToResponse(),
		Severity: condition.Severity,
	}
}

func (events ConditionEvents) ToResponse() ConditionEventsResponse {
	var resp = []ConditionEventMetaResponse{}
	for _, s := range events {
		resp = append(resp,
			ConditionEventMetaResponse{
				ID:   s.ID,
				Date: s.Date,
			})
	}
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Date.Sub(resp[j].Date) > 0
	})
	return ConditionEventsResponse{ConditionEvents: resp}
}

func (event ConditionEvent) ToResponse() ConditionEventResponse {
	var conditionResp = []ConditionResponse{}
	for _, con := range event.Conditions {
		conditionResp = append(conditionResp, con.ToResponse())
	}
	var resp = ConditionEventResponse{
		ID:         event.ID,
		Date:       event.Date,
		Conditions: conditionResp,
	}
	return resp
}

type PostConditionResponse struct {
	Condition ConditionResponse         `json:"condition"`
	Symptoms  SymptomCategoriesResponse `json:"symptoms"`
}
