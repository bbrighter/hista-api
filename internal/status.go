package internal

import (
	"errors"
	"time"

	"encore.app/entity"
)

type IStatusRepo interface {
	Find() entity.Statuses
	First(*entity.Status) error
	Save(*entity.Status) error
	Delete(*entity.Status) error
	FindForDate(time.Time) (entity.Status, bool)
}

type IStatusUseCase interface {
	Find() entity.Statuses
	Create(time.Time) (entity.Status, error)
	SaveMorning(statusId uint, date time.Time, morningStatusId uint, fitness entity.Quality, sleep entity.Quality) (entity.Status, error)
	SaveEvening(statusId uint, date time.Time, eveningStatusId uint, fitness entity.Quality) (entity.Status, error)
	Delete(id uint) error
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
	uc.repo.Save(&status)
	return status, nil
}

func (uc StatusUseCase) SaveMorning(statusId uint, date time.Time, morningStatusId uint, fitness entity.Quality, sleep entity.Quality) (entity.Status, error) {
	status := entity.Status{ID: statusId, Date: date,
		Morning: &entity.MorningStatus{Fitness: fitness, Sleep: sleep, StatusID: statusId, ID: morningStatusId},
	}
	err := uc.repo.Save(&status)
	if err != nil {
		return status, err
	}
	err = uc.repo.First(&status)
	return status, err
}

func (uc StatusUseCase) SaveEvening(statusId uint, date time.Time, eveningStatusId uint, fitness entity.Quality) (entity.Status, error) {
	status := entity.Status{ID: statusId, Date: date,
		Evening: &entity.EveningStatus{Fitness: fitness, StatusID: statusId, ID: eveningStatusId},
	}
	err := uc.repo.Save(&status)
	if err != nil {
		return status, err
	}
	err = uc.repo.First(&status)
	return status, err
}

func (uc StatusUseCase) Delete(id uint) error {
	var status = entity.Status{ID: id}
	return uc.repo.Delete(&status)
}
