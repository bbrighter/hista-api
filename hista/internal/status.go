package internal

import (
	"context"
	"time"

	"encore.app/errors"
	"encore.app/hista/entity"
	"encore.dev/beta/errs"
)

type IStatusRepo interface {
	Find(ctx context.Context) ([]*entity.Status, error)
	Create(ctx context.Context, status *entity.Status) error
	Update(ctx context.Context, status *entity.Status) error
	Delete(ctx context.Context, id uint) error
	FindForDate(ctx context.Context, date time.Time) (entity.Status, bool)
}

type IStatusUseCase interface {
	Find(ctx context.Context) (entity.Statuses, error)
	Create(ctx context.Context, date time.Time) (entity.Status, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, statusId uint, date time.Time, morning entity.MorningStatus, evening entity.EveningStatus) error
}

type StatusUseCase struct {
	repo IStatusRepo
}

func NewStatusUseCase(repo IStatusRepo) StatusUseCase {
	return StatusUseCase{repo: repo}
}

func (uc StatusUseCase) Find(ctx context.Context) (entity.Statuses, error) {
	statuses, err := uc.repo.Find(ctx)
	return statuses, errors.MapError(err)
}

func (uc StatusUseCase) Create(ctx context.Context, date time.Time) (entity.Status, error) {
	_, exists := uc.repo.FindForDate(ctx, date)
	if exists {
		return entity.Status{}, &errs.Error{Code: errs.AlreadyExists, Message: "already exists"}
	}
	status := entity.Status{Date: date}
	err := uc.repo.Create(ctx, &status)
	return status, errors.MapError(err)
}

func (uc StatusUseCase) Update(
	ctx context.Context,
	statusId uint,
	date time.Time,
	morning entity.MorningStatus,
	evening entity.EveningStatus,
) error {
	status := &entity.Status{
		ID:             statusId,
		Date:           date,
		MorningFitness: morning.Fitness,
		EveningFitness: evening.Fitness,
		MorningSleep:   morning.Sleep,
	}
	return errors.MapError(uc.repo.Update(ctx, status))

}

func (uc StatusUseCase) Delete(ctx context.Context, id uint) error {
	return errors.MapError(uc.repo.Delete(ctx, id))
}
