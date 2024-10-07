package api

import (
	"context"
	"time"

	"encore.app/entity"
)

type StatusParams struct {
	Date      time.Time        `json:"date"`
	TimeOfDay entity.TimeOfDay `json:"timeOfDay"`
	Fitness   entity.Quality   `json:"fitness"`
	Sleep     *entity.Quality  `json:"sleep,omitempty"`
}

// encore:api auth method=POST path=/status
func (service *Service) CreateStatus(ctx context.Context, params StatusParams) (entity.StatusResponse, error) {
	status, err := service.status.Create(params.Date, params.TimeOfDay, params.Fitness, params.Sleep)
	return status.ToResp(), err
}

// encore:api auth method=GET path=/status
func (service *Service) ListStatus(ctx context.Context) (entity.StatusesResponse, error) {
	statuses := service.status.Find()
	return statuses.ToResp(), nil
}
