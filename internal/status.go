package internal

import (
	"errors"
	"time"

	"encore.app/entity"
)

type IStatusRepo interface {
	Find() entity.Statuses
	First(*entity.Status) error
	Create(*entity.Status) error
	Update(*entity.Status) error
	Delete(*entity.Status) error
	FindForDate(time.Time) (entity.Status, bool)
}

type IStatusUseCase interface {
	Find() entity.Statuses
	Create(time.Time) (entity.Status, error)
	Delete(id uint) error
	Update(statusId uint, date time.Time, morning entity.MorningStatus, evening entity.EveningStatus) error
}

type StatusUseCase struct {
	repo IStatusRepo
}

func NewStatusUseCase(repo IStatusRepo) StatusUseCase {
	return StatusUseCase{repo: repo}
}

func (uc StatusUseCase) Find() entity.Statuses {
	return uc.repo.Find()
}

func (uc StatusUseCase) Create(date time.Time) (entity.Status, error) {
	_, exists := uc.repo.FindForDate(date)
	if exists {
		return entity.Status{}, errors.New("status for date already exists")
	}
	status := entity.Status{Date: date}
	uc.repo.Create(&status)
	return status, nil
}

func (uc StatusUseCase) Update(
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
	return uc.repo.Update(status)

}

func (uc StatusUseCase) Delete(id uint) error {
	var status = entity.Status{ID: id}
	return uc.repo.Delete(&status)
}
