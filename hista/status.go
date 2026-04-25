package hista

import (
	"context"
	"slices"
	"time"

	"encore.app/errors"
	"encore.app/hista/internal/status"
	"encore.dev/types/option"
	"encore.dev/types/uuid"
)

type StatusResponse struct {
	ID             uint               `json:"id"`
	Date           time.Time          `json:"date"`
	MorningFitness option.Option[int] `json:"morningFitness" encore:"optional"`
	EveningFitness option.Option[int] `json:"eveningFitness" encore:"optional"`
	MorningSleep   option.Option[int] `json:"morningSleep" encore:"optional"`
}

type StatusListResponse struct {
	Statuses []StatusResponse `json:"statuses"`
}

func toStatusListResp(statuses []*status.Status) StatusListResponse {
	responses := []StatusResponse{}
	for _, s := range statuses {
		responses = append(responses, toStatusResp(s))
	}
	slices.SortFunc(responses, func(a, b StatusResponse) int {
		return b.Date.Compare(a.Date)
	})
	return StatusListResponse{Statuses: responses}
}

func toStatusResp(s *status.Status) StatusResponse {
	return StatusResponse{
		ID:             s.ID,
		Date:           s.Date,
		MorningFitness: option.FromPointer(s.MorningFitness),
		EveningFitness: option.FromPointer(s.EveningFitness),
		MorningSleep:   option.FromPointer(s.MorningSleep),
	}
}

type DateParam struct {
	Date time.Time `json:"date"`
}

// encore:api auth method=POST path=/piid/:piid/status
func (service *Service) PostStatus(ctx context.Context, piid uuid.UUID, params DateParam) (StatusResponse, error) {
	status, err := service.status.CreateStatus(ctx, params.Date)
	if err != nil {
		return StatusResponse{}, errors.MapError(err)
	}
	return toStatusResp(status), nil
}

type PatchStatusParams struct {
	Date           option.Option[time.Time] `json:"date"`
	MorningFitness option.Option[int]       `json:"morningFitness" encore:"optional"`
	MorningSleep   option.Option[int]       `json:"morningSleep" encore:"optional"`
	EveningFitness option.Option[int]       `json:"eveningFitness" encore:"optional"`
}

// encore:api auth method=PATCH path=/piid/:piid/status/:id
func (service *Service) PatchStatus(ctx context.Context, piid uuid.UUID, id uint, params PatchStatusParams) error {
	err := service.status.UpdateStatus(ctx, id, struct {
		Date           *time.Time
		MorningFitness *int
		MorningSleep   *int
		EveningFitness *int
	}{
		Date:           params.Date.PtrOrNil(),
		MorningFitness: params.MorningFitness.PtrOrNil(),
		MorningSleep:   params.MorningSleep.PtrOrNil(),
		EveningFitness: params.EveningFitness.PtrOrNil(),
	})
	return errors.MapError(err)
}

// encore:api auth method=GET path=/piid/:piid/status
func (service *Service) ListStatus(ctx context.Context, piid uuid.UUID) (StatusListResponse, error) {
	statuses, err := service.status.ListStatuses(ctx)
	if err != nil {
		return StatusListResponse{}, errors.MapError(err)
	}
	return toStatusListResp(statuses), nil
}

// encore:api auth method=DELETE path=/piid/:piid/status/:id
func (service *Service) DeleteStatus(ctx context.Context, piid uuid.UUID, id uint) error {
	err := service.status.DeleteStatus(ctx, id)
	return errors.MapError(err)
}
