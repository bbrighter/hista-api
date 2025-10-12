package internal

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/entity"
	"github.com/stretchr/testify/assert"
)

type HeadacheTestRepo struct{}

func (r HeadacheTestRepo) ListHeadaches(ctx context.Context) ([]*entity.Headache, error) {
	return nil, nil
}
func (r HeadacheTestRepo) CreateHeadache(ctx context.Context, ha *entity.Headache) error {
	ha.ID = 10
	return nil
}
func (r HeadacheTestRepo) DeleteHeadache(ctx context.Context, haId uint) error {
	return nil
}
func (r HeadacheTestRepo) GetHeadache(ctx context.Context, haId uint) (*entity.Headache, error) {
	return &entity.Headache{}, nil
}
func (r HeadacheTestRepo) PatchHeadache(ctx context.Context, haId uint, date *time.Time, severity *entity.HeadacheSeverity, types *entity.HeadacheTypes, positions *entity.HeadachePositions, symptoms *entity.HeadacheSymptoms, description *string) error {
	return nil
}

func TestCreateHeadache(t *testing.T) {
	ctx := t.Context()
	uc := NewHeadacheUseCase(HeadacheTestRepo{})
	date := time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC)
	var severity entity.HeadacheSeverity = 2
	id, err := uc.Create(ctx, date, severity)
	assert.EqualValues(t, 10, id)
	assert.NoError(t, err)

}
