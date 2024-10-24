package entity

import (
	"slices"
	"time"
)

type Status struct {
	ID      uint
	Date    time.Time
	Morning *MorningStatus `gorm:"constraint:OnDelete:CASCADE"`
	Evening *EveningStatus `gorm:"constraint:OnDelete:CASCADE"`
}

type MorningStatus struct {
	ID       uint    `json:"id"`
	StatusID uint    `json:"statusId" gorm:"unique"`
	Fitness  Quality `json:"fitness"`
	Sleep    Quality `json:"sleep"`
}

type EveningStatus struct {
	ID       uint    `json:"id"`
	StatusID uint    `json:"statusId" gorm:"unique"`
	Fitness  Quality `json:"fitness"`
}

type Statuses []Status

type StatusesResponse struct {
	Statuses []StatusResponse `json:"statuses"`
}

func (statuses Statuses) ToResp() StatusesResponse {
	responses := []StatusResponse{}
	for _, s := range statuses {
		responses = append(responses, StatusResponse(s))
	}
	slices.SortFunc(responses, func(a, b StatusResponse) int {
		return b.Date.Compare(a.Date)
	})
	return StatusesResponse{Statuses: responses}
}

type StatusResponse struct {
	ID      uint           `json:"id"`
	Date    time.Time      `json:"date"`
	Morning *MorningStatus `json:"morning,omitempty" encore:"optional"`
	Evening *EveningStatus `json:"evening,omitempty" encore:"optional"`
}

func (status Status) ToResp() StatusResponse {
	return StatusResponse(status)
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
