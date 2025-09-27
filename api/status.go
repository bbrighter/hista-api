package api

import (
	"context"
	"time"

	"encore.app/entity"
)

type DateParam struct {
	Date time.Time `json:"date"`
}

// encore:api auth method=POST path=/status
func (service *Service) PostStatus(ctx context.Context, params DateParam) (entity.StatusResponse, error) {
	status, err := service.status.Create(ctx, params.Date)
	return status.ToResp(), err
}

type PatchStatusParams struct {
	Date           time.Time `json:"date" encore:"optional"`
	MorningFitness *int      `json:"morningFitness" encore:"optional"`
	MorningSleep   *int      `json:"morningSleep" encore:"optional"`
	EveningFitness *int      `json:"eveningFitness" encore:"optional"`
}

// encore:api auth method=PATCH path=/status/:id
func (service *Service) PatchStatus(ctx context.Context, id uint, params PatchStatusParams) error {
	morningStatus := entity.MorningStatus{Fitness: params.MorningFitness, Sleep: params.MorningSleep}
	eveningStatus := entity.EveningStatus{Fitness: params.EveningFitness}
	err := service.status.Update(ctx, id, params.Date, morningStatus, eveningStatus)
	return err
}

// encore:api auth method=GET path=/status
func (service *Service) ListStatus(ctx context.Context) (entity.StatusesResponse, error) {
	statuses, err := service.status.Find(ctx)
	return statuses.ToResp(), err
}

// encore:api auth method=DELETE path=/status/:id
func (service *Service) DeleteStatus(ctx context.Context, id uint) error {
	return service.status.Delete(ctx, id)
}
