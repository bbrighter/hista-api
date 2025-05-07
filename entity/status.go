package entity

import (
	"slices"
	"time"
)

type Status struct {
	ID             uint
	Date           time.Time
	MorningFitness *int
	EveningFitness *int
	MorningSleep   *int
}

type StatusResponse struct {
	ID             uint      `json:"id"`
	Date           time.Time `json:"date"`
	MorningFitness *int      `json:"morningFitness" encore:"optional"`
	EveningFitness *int      `json:"eveningFitness" encore:"optional"`
	MorningSleep   *int      `json:"morningSleep" encore:"optional"`
}

type Statuses []Status

type StatusesResponse struct {
	Statuses []StatusResponse `json:"statuses"`
}

func (statuses Statuses) ToResp() StatusesResponse {
	responses := []StatusResponse{}
	for _, s := range statuses {
		responses = append(responses, s.ToResp())
	}
	slices.SortFunc(responses, func(a, b StatusResponse) int {
		return b.Date.Compare(a.Date)
	})
	return StatusesResponse{Statuses: responses}
}

type MorningStatus struct {
	Fitness *int `json:"fitness"`
	Sleep   *int `json:"sleep"`
}

type EveningStatus struct {
	Fitness *int `json:"fitness"`
}

func (s Status) ToResp() StatusResponse {
	return StatusResponse(s)
}
