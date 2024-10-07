package internal

import (
	"time"

	"encore.app/entity"
)

type StatusUseCase struct {
	repo IStatusRepo
}

func NewStatusUseCase(repo IStatusRepo) StatusUseCase {
	return StatusUseCase{repo: repo}
}

func (uc StatusUseCase) Find() entity.Statuses {
	return uc.repo.Find()
}

func (uc StatusUseCase) Create(
	date time.Time,
	timeOfDay entity.TimeOfDay,
	fitness entity.Quality,
	sleep *entity.Quality,
) (entity.Status, error) {
	var status = entity.Status{
		Date:      date,
		TimeOfDay: timeOfDay,
		Fitness:   fitness,
		Sleep:     sleep,
	}
	validator, err := status.GetValidator()
	if err != nil {
		return status, err
	}
	err = validator.Validate(sleep)
	if err != nil {
		return status, err
	}
	err = uc.repo.Create(&status)
	return status, err
}
