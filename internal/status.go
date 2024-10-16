package internal

import (
	"errors"
	"time"

	"encore.app/entity"
)

type IStatusRepo interface {
	Find() entity.Statuses
	Save(*entity.Status) error
	Delete(*entity.Status) error
	FindForDate(time.Time) (entity.Status, bool)
}

type IStatusUseCase interface {
	Find() entity.Statuses
	CreateMorning(date time.Time, fitness entity.Quality, sleep entity.Quality) (entity.Status, error)
	CreateEvening(date time.Time, fitness entity.Quality) (entity.Status, error)
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

func (uc StatusUseCase) CreateMorning(date time.Time, fitness entity.Quality, sleep entity.Quality) (entity.Status, error) {
	status, exists := uc.repo.FindForDate(date)
	if !exists {
		status = entity.Status{Date: date}
	}
	if status.Morning != nil {
		return status, errors.New("morning status already exists")
	}
	status.Morning = &entity.MorningStatus{
		Fitness: fitness,
		Sleep:   sleep,
	}
	err := uc.repo.Save(&status)
	return status, err
}

func (uc StatusUseCase) CreateEvening(date time.Time, fitness entity.Quality) (entity.Status, error) {
	status, exists := uc.repo.FindForDate(date)
	if !exists {
		status = entity.Status{Date: date}
	}
	if status.Evening != nil {
		return status, errors.New("evening status already exists")
	}
	status.Evening = &entity.EveningStatus{
		Fitness: fitness,
	}
	err := uc.repo.Save(&status)
	return status, err
}

func (uc StatusUseCase) Delete(id uint) error {
	var status = entity.Status{ID: id}
	return uc.repo.Delete(&status)
}
