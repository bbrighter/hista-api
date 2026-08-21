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
	ID                    uint               `json:"id"`
	Date                  time.Time          `json:"date"`
	MorningFitness        option.Option[int] `json:"morningFitness" `
	DayFitness            option.Option[int] `json:"dayFitness"`
	EveningFitness        option.Option[int] `json:"eveningFitness" `
	MorningSleep          option.Option[int] `json:"morningSleep" `
	Depressive            option.Option[int] `json:"depressive" `
	Tense                 option.Option[int] `json:"tense" `
	MoodSwings            option.Option[int] `json:"moodSwings" `
	Irritable             option.Option[int] `json:"irritable" `
	LossOfInterest        option.Option[int] `json:"lossOfInterest" `
	ConcentrationProblems option.Option[int] `json:"concentrationProblems" `
	LackOfDrive           option.Option[int] `json:"lackOfDrive" `
	AppetiteChanges       option.Option[int] `json:"appetiteChanges" `
	SleepProblems         option.Option[int] `json:"sleepProblems" `
	Overwhelmed           option.Option[int] `json:"overwhelmed" `
	Crash                 bool               `json:"crash"`
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
		ID:                    s.ID,
		Date:                  s.Date,
		MorningFitness:        option.FromPointer(s.MorningFitness),
		DayFitness:            option.FromPointer(s.DayFitness),
		EveningFitness:        option.FromPointer(s.EveningFitness),
		MorningSleep:          option.FromPointer(s.MorningSleep),
		Depressive:            option.FromPointer(s.Depressive),
		Tense:                 option.FromPointer(s.Tense),
		MoodSwings:            option.FromPointer(s.MoodSwings),
		Irritable:             option.FromPointer(s.Irritable),
		LossOfInterest:        option.FromPointer(s.LossOfInterest),
		ConcentrationProblems: option.FromPointer(s.ConcentrationProblems),
		LackOfDrive:           option.FromPointer(s.LackOfDrive),
		AppetiteChanges:       option.FromPointer(s.AppetiteChanges),
		SleepProblems:         option.FromPointer(s.SleepProblems),
		Overwhelmed:           option.FromPointer(s.Overwhelmed),
		Crash:                 s.Crash,
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
	Date                  option.Option[time.Time] `json:"date"`
	MorningFitness        option.Option[int]       `json:"morningFitness" `
	MorningSleep          option.Option[int]       `json:"morningSleep" `
	DayFitness            option.Option[int]       `json:"dayFitness"`
	EveningFitness        option.Option[int]       `json:"eveningFitness" `
	Depressive            option.Option[int]       `json:"depressive" `
	Tense                 option.Option[int]       `json:"tense" `
	MoodSwings            option.Option[int]       `json:"moodSwings" `
	Irritable             option.Option[int]       `json:"irritable" `
	LossOfInterest        option.Option[int]       `json:"lossOfInterest" `
	ConcentrationProblems option.Option[int]       `json:"concentrationProblems" `
	LackOfDrive           option.Option[int]       `json:"lackOfDrive" `
	AppetiteChanges       option.Option[int]       `json:"appetiteChanges" `
	SleepProblems         option.Option[int]       `json:"sleepProblems" `
	Overwhelmed           option.Option[int]       `json:"overwhelmed" `
	Crash                 bool                     `json:"crash"`
}

// encore:api auth method=PATCH path=/piid/:piid/status/:id
func (service *Service) PatchStatus(ctx context.Context, piid uuid.UUID, id uint, params PatchStatusParams) error {
	err := service.status.UpdateStatus(ctx, id, status.UpdateStatusParams{
		Date:                  params.Date.PtrOrNil(),
		MorningFitness:        params.MorningFitness.PtrOrNil(),
		MorningSleep:          params.MorningSleep.PtrOrNil(),
		DayFitness:            params.DayFitness.PtrOrNil(),
		EveningFitness:        params.EveningFitness.PtrOrNil(),
		Depressive:            params.Depressive.PtrOrNil(),
		Tense:                 params.Tense.PtrOrNil(),
		MoodSwings:            params.MoodSwings.PtrOrNil(),
		Irritable:             params.Irritable.PtrOrNil(),
		LossOfInterest:        params.LossOfInterest.PtrOrNil(),
		ConcentrationProblems: params.ConcentrationProblems.PtrOrNil(),
		LackOfDrive:           params.LackOfDrive.PtrOrNil(),
		AppetiteChanges:       params.AppetiteChanges.PtrOrNil(),
		SleepProblems:         params.SleepProblems.PtrOrNil(),
		Overwhelmed:           params.Overwhelmed.PtrOrNil(),
		Crash:                 &params.Crash,
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
