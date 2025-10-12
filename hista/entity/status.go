package entity

import (
	"slices"
	"time"

	"encore.dev/types/uuid"
)

type Status struct {
	ID             uint      `gorm:"primaryKey"`
	PIID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Date           time.Time
	MorningFitness *int
	EveningFitness *int
	MorningSleep   *int
}

func (s *Status) SetPiid(id uuid.UUID) {
	s.PIID = id
}

type StatusResponse struct {
	ID             uint      `json:"id"`
	Date           time.Time `json:"date"`
	MorningFitness *int      `json:"morningFitness" encore:"optional"`
	EveningFitness *int      `json:"eveningFitness" encore:"optional"`
	MorningSleep   *int      `json:"morningSleep" encore:"optional"`
}

type Statuses []*Status

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
	return StatusResponse{
		ID:             s.ID,
		Date:           s.Date,
		MorningFitness: s.MorningFitness,
		EveningFitness: s.EveningFitness,
		MorningSleep:   s.MorningSleep,
	}
}
