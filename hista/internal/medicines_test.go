package internal

import (
	"context"
	"testing"

	"encore.app/hista/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMedicineRepo struct {
	mock.Mock
}

func (m *MockMedicineRepo) List(ctx context.Context) ([]*entity.Medicine, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Medicine), args.Error(1)
}

func (m *MockMedicineRepo) Create(ctx context.Context, name string) (uint, error) {
	args := m.Called(ctx, name)
	return uint(args.Int(0)), args.Error(1)
}

func (m *MockMedicineRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockMedicineRepo) Patch(ctx context.Context, id uint, column string, value any) error {
	return m.Called(ctx, id, column, value).Error(0)
}

func (m *MockMedicineRepo) BulkShift(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *MockMedicineRepo) Reorder(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

type MockUnitOfWork struct {
	mock.Mock
	repo *MockMedicineRepo
}

func (m *MockUnitOfWork) WithTransaction(
	ctx context.Context,
	fn func(tx UnitOfWork) error,
) error {
	return fn(m)
}

func (m *MockUnitOfWork) Medicine() IMedicineRepo {
	return m.repo
}

func TestMedicineMgmtUseCase_Reorder(t *testing.T) {
	ctx := context.Background()

	var id1 uint = 1
	var id2 uint = 2

	tests := map[string]struct {
		initialList   []*entity.Medicine
		afterRecalc   []*entity.Medicine
		prevId        *uint
		nextId        *uint
		expectRecalc  bool
		expectedOrder int
	}{
		"insert between": {
			initialList: []*entity.Medicine{
				{ID: 1, SortOrder: 100},
				{ID: 2, SortOrder: 200},
			},
			prevId:        &id1,
			nextId:        &id2,
			expectRecalc:  false,
			expectedOrder: 150,
		},
		"needs recalc": {
			initialList: []*entity.Medicine{
				{ID: 1, SortOrder: 100},
				{ID: 2, SortOrder: 101},
			},
			prevId: &id1,
			nextId: &id2,
			afterRecalc: []*entity.Medicine{
				{ID: 1, SortOrder: 100},
				{ID: 2, SortOrder: 200},
			},
			expectRecalc:  true,
			expectedOrder: 150,
		},
		"no prev id, no resorting": {
			initialList: []*entity.Medicine{
				{ID: 1, SortOrder: 100},
				{ID: 2, SortOrder: 200},
			},
			nextId:        &id1,
			expectRecalc:  false,
			expectedOrder: 50,
		},
		"no prev id, resorting": {
			initialList: []*entity.Medicine{
				{ID: 1, SortOrder: 1},
				{ID: 2, SortOrder: 200},
			},
			afterRecalc: []*entity.Medicine{
				{ID: 1, SortOrder: 100},
				{ID: 2, SortOrder: 200},
			},
			nextId:        &id1,
			expectRecalc:  true,
			expectedOrder: 50,
		},
		"no next id, no resorting": {
			initialList: []*entity.Medicine{
				{ID: 1, SortOrder: 100},
				{ID: 2, SortOrder: 200},
			},
			prevId:        &id2,
			expectRecalc:  false,
			expectedOrder: 300,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			repo := new(MockMedicineRepo)
			uow := &MockUnitOfWork{repo: repo}
			usecase := NewMedicineMgtmUseCase(repo, uow)

			repo.On("List", ctx).Return(test.initialList, nil).Once()

			if test.expectRecalc {
				repo.On("Reorder", ctx).Return(nil).Once()
				repo.On("List", ctx).Return(test.afterRecalc, nil).Once()
			}

			repo.On("Patch", ctx, uint(3), "sort_order", test.expectedOrder).Return(nil).Once()

			err := usecase.Reorder(ctx, 3, test.prevId, test.nextId)

			assert.NoError(t, err)
			repo.AssertExpectations(t)
		})
	}
}

func TestMedicineMgmtUseCase_Create(t *testing.T) {
	ctx := context.Background()

	repo := new(MockMedicineRepo)
	uow := &MockUnitOfWork{repo: repo}
	usecase := NewMedicineMgtmUseCase(repo, uow)

	repo.On("BulkShift", ctx).Return(nil).Once()
	repo.On("Create", ctx, "name").Return(1, nil).Once()

	id, err := usecase.Create(ctx, "name")

	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)
	repo.AssertExpectations(t)
}
