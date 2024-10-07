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

func (repo TestStatusRepo) Create(status *entity.Status) error {
	status.ID = 1
	return nil
}

func TestCreate(t *testing.T) {
	repo := TestStatusRepo{}
	uc := internal.NewStatusUseCase(repo)

	var status entity.Status
	var err error
	status, err = uc.Create(time.Now(), entity.Evening, entity.Bad, nil)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, status.ID)

	sleep := entity.Bad
	_, err = uc.Create(time.Now(), entity.Evening, entity.Bad, &sleep)
	assert.Error(t, err)

	_, err = uc.Create(time.Now(), entity.Morning, entity.Bad, nil)
	assert.Error(t, err)

	_, err = uc.Create(time.Now(), entity.Morning, entity.Bad, &sleep)
	assert.NoError(t, err)
}
