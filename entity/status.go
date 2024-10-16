package entity

import (
	"time"
)

type Status struct {
	ID      uint           `json:"id"`
	Date    time.Time      `json:"date"`
	Morning *MorningStatus `json:"morning,omitempty" encore:"optional" gorm:"constraint:OnDelete:CASCADE"`
	Evening *EveningStatus `json:"evening,omitempty" encore:"optional" gorm:"constraint:OnDelete:CASCADE"`
}

type MorningStatus struct {
	ID       uint    `json:"id"`
	StatusID uint    `json:"statusId"`
	Fitness  Quality `json:"fitness"`
	Sleep    Quality `json:"sleep"`
}

type EveningStatus struct {
	ID       uint    `json:"id"`
	StatusID uint    `json:"statusId"`
	Fitness  Quality `json:"fitness"`
}

type Statuses []Status

type StatusesResponse struct {
	Statuses Statuses `json:"statuses"`
}

func (statuses Statuses) ToResp() StatusesResponse {
	return StatusesResponse{Statuses: statuses}
}

type TimeOfDay string

const (
	Morning TimeOfDay = "morning"
	Evening TimeOfDay = "evening"
)

func (tod TimeOfDay) IsValid() bool {
	switch tod {
	case Morning, Evening:
		return true
	}
	return false
}
