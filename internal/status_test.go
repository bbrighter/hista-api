package internal_test

import (
	"testing"
	"time"

	"encore.app/entity"
	"encore.app/internal"
	"github.com/stretchr/testify/assert"
)

type TestStatusRepo struct{}

func (repo TestStatusRepo) Find() entity.Statuses {
	return entity.Statuses{}
}

func (repo TestStatusRepo) Save(status *entity.Status) error {
	status.ID = 1
	return nil
}

func (repo TestStatusRepo) Delete(status *entity.Status) error {
	status.ID = 0
	return nil
}

func (repo TestStatusRepo) FindForDate(date time.Time) (entity.Status, bool) {
	status := entity.Status{ID: 1, Date: time.Now(), Morning: &entity.MorningStatus{ID: 1}}
	return status, true
}

func TestCreate(t *testing.T) {
	repo := TestStatusRepo{}
	uc := internal.NewStatusUseCase(repo)

	var status entity.Status
	var err error
	status, err = uc.CreateMorning(time.Now(), entity.Bad, entity.Good)
	assert.EqualError(t, err, "morning status already exists")

	status, err = uc.CreateEvening(time.Now(), entity.Bad)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, status.ID)
}
