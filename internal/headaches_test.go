package internal

import (
	"testing"
	"time"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
)

type HeadacheTestRepo struct{}

func (r HeadacheTestRepo) ListHeadaches() entity.Headaches {
	return entity.Headaches{}
}
func (r HeadacheTestRepo) CreateHeadache(ha *entity.Headache) error {
	ha.ID = 10
	return nil
}
func (r HeadacheTestRepo) DeleteHeadache(haId uint) error {
	return nil
}
func (r HeadacheTestRepo) GetHeadache(haId uint) (entity.Headache, error) {
	return entity.Headache{}, nil
}
func (r HeadacheTestRepo) PatchHeadache(haId uint, date *time.Time, severity *entity.HeadacheSeverity, types *entity.HeadacheTypes, positions *entity.HeadachePositions, symptoms *entity.HeadacheSymptoms) error {
	return nil
}

func TestCreateHeadache(t *testing.T) {
	uc := NewHeadacheUseCase(HeadacheTestRepo{})
	date := time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC)
	var severity entity.HeadacheSeverity = 2
	id, err := uc.Create(date, severity)
	assert.EqualValues(t, 10, id)
	assert.NoError(t, err)

}
