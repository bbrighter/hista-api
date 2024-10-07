package entity

import (
	"fmt"
	"time"
)

type Status struct {
	ID        uint
	Date      time.Time
	TimeOfDay TimeOfDay
	Fitness   Quality
	Sleep     *Quality
}

type Statuses []Status

type TimeOfDay string

const (
	Morning TimeOfDay = "Morning"
	Evening TimeOfDay = "Evening"
)

type Validator interface {
	Validate(sleep *Quality) error
}

type MorningValidator struct{}

func (ms MorningValidator) Validate(sleep *Quality) error {
	if sleep == nil {
		return fmt.Errorf("sleep must not be nil")
	}
	return nil
}

type EveningValidator struct{}

func (es EveningValidator) Validate(sleep *Quality) error {
	if sleep != nil {
		return fmt.Errorf("sleep must be nil")
	}
	return nil
}

func (status Status) GetValidator() (Validator, error) {
	switch status.TimeOfDay {
	case Morning:
		return MorningValidator{}, nil
	case Evening:
		return EveningValidator{}, nil
	default:
		return nil, fmt.Errorf("invalid time of day: %v", status.TimeOfDay)
	}
}

type StatusResponse struct {
	ID        uint      `json:"id"`
	Date      time.Time `json:"date"`
	TimeOfDay TimeOfDay `json:"timeOfDay"`
	Fitness   Quality   `json:"fitness"`
	Sleep     *Quality  `json:"sleep,omitempty"`
}

func (status Status) ToResp() StatusResponse {
	return StatusResponse(status)
}

type StatusesResponse struct {
	Statuses []StatusResponse `json:"statuses"`
}

func (statuses Statuses) ToResp() StatusesResponse {
	var stats = []StatusResponse{}
	for _, status := range statuses {
		stats = append(stats, status.ToResp())
	}
	return StatusesResponse{Statuses: stats}
}
