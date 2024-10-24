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
	ID        uint             `json:"id,omitempty" encore:"optional"`
}

type DateParam struct {
	Date time.Time `json:"date"`
}

// encore:api auth method=POST path=/status
func (service *Service) PostStatus(ctx context.Context, params DateParam) (entity.StatusResponse, error) {
	status, err := service.status.Create(params.Date)
	return status.ToResp(), err
}

// encore:api auth method=PUT path=/status/:id
func (service *Service) PutStatus(ctx context.Context, id uint, params StatusParams) (entity.StatusResponse, error) {
	switch params.TimeOfDay {
	case entity.Morning:
		status, err := service.status.SaveMorning(id, params.Date, params.ID, params.Fitness, params.Sleep)
		return status.ToResp(), err
	case entity.Evening:
		status, err := service.status.SaveEvening(id, params.Date, params.ID, params.Fitness)
		return status.ToResp(), err
	default:
		return entity.StatusResponse{}, errors.BadRequestf("timeOfDay must be a valid value %v", params.TimeOfDay)
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
