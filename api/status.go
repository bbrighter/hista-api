package api

import (
	"context"
	"time"

	"encore.app/entity"
	"encore.app/errors"
)

type StatusParams struct {
	TimeOfDay entity.TimeOfDay `json:"timeOfDay"`
	Date      time.Time        `json:"date"`
	Fitness   entity.Quality   `json:"fitness"`
	Sleep     entity.Quality   `json:"sleep,omitempty" encore:"optional"`
}

// encore:api auth method=POST path=/status
func (service *Service) CreateStatus(ctx context.Context, params StatusParams) (entity.IDResponse, error) {
	switch params.TimeOfDay {
	case entity.Morning:
		status, err := service.status.CreateMorning(params.Date, params.Fitness, params.Sleep)
		return entity.IDResponse{ID: status.ID}, err
	case entity.Evening:
		status, err := service.status.CreateEvening(params.Date, params.Fitness)
		return entity.IDResponse{ID: status.ID}, err
	default:
		return entity.IDResponse{}, errors.BadRequestf("timeOfDay must be valid value %v", params.TimeOfDay)
	}
}

// encore:api auth method=GET path=/status
func (service *Service) ListStatus(ctx context.Context) (entity.StatusesResponse, error) {
	statuses := service.status.Find()
	return statuses.ToResp(), nil
}

// encore:api auth method=DELETE path=/status/:id
func (service *Service) DeleteStatus(ctx context.Context, id uint) error {
	return service.status.Delete(id)
}
