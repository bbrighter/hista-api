package api

import (
	"context"
	"time"

	"encore.app/entity"
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

type PatchStatusParams struct {
	Date    time.Time      `json:"date"`
	Morning *MorningParams `json:"morning" encore:"optional"`
	Evening *EveningParams `json:"evening" encore:"optional"`
}

type MorningParams struct {
	Fitness entity.Quality `json:"fitness"`
	Sleep   entity.Quality `json:"sleep,omitempty" encore:"optional"`
	ID      uint           `json:"id,omitempty" encore:"optional"`
}

type EveningParams struct {
	Fitness entity.Quality `json:"fitness"`
	ID      uint           `json:"id,omitempty" encore:"optional"`
}

// encore:api auth method=PUT path=/status/:id
func (service *Service) PutStatus(ctx context.Context, id uint, params PatchStatusParams) (entity.StatusResponse, error) {
	morning := new(entity.MorningStatus)
	if params.Morning != nil {
		morning.Fitness = params.Morning.Fitness
		morning.Sleep = params.Morning.Sleep
		morning.StatusID = id
		if params.Morning.ID > 0 {
			morning.StatusID = params.Morning.ID
		}
	}
	evening := new(entity.EveningStatus)
	if params.Evening != nil {
		evening.Fitness = params.Evening.Fitness
		evening.StatusID = id
		if params.Evening.ID > 0 {
			evening.StatusID = params.Evening.ID
		}
	}
	status, err := service.status.Update(id, params.Date, morning, evening)
	return status.ToResp(), err
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
