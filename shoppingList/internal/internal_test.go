package internal

import (
	"context"
	"testing"

	"encore.app/shoppingList/entity"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockItemRepo struct{ mock.Mock }
type MockUow struct{ mock.Mock }

func (m *MockItemRepo) Create(ctx context.Context, productId uint, listId uint) (uint, error) {
	args := m.Called(ctx, productId, listId)
	return uint(args.Int(0)), args.Error(1)
}

func (m *MockItemRepo) Delete(ctx context.Context, itemId uint) error {
	args := m.Called(ctx, itemId)
	return args.Error(0)
}

func (m *MockItemRepo) Check(ctx context.Context, itemId uint) error {
	args := m.Called(ctx, itemId)
	return args.Error(0)
}

func (m *MockItemRepo) Find(ctx context.Context, id uint) (*entity.Item, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Item), args.Error(1)
}
func (m *MockItemRepo) List(ctx context.Context, listId uint) ([]entity.Item, error) {
	args := m.Called(ctx, listId)
	return args.Get(0).([]entity.Item), args.Error(1)
}

func (m *MockItemRepo) CheckUniqueness(ctx context.Context, productId uint, listId uint) error {
	args := m.Called(ctx, productId, listId)
	return args.Error(0)
}

type MockListRepo struct{ mock.Mock }

func (m *MockListRepo) First(ctx context.Context) (*entity.List, error) {
	args := m.Called(ctx)
	return args.Get(0).(*entity.List), args.Error(1)
}
func (m *MockListRepo) Create(ctx context.Context) (uint, error) {
	args := m.Called(ctx)
	return uint(args.Int(0)), args.Error(1)
}
func (m *MockListRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockProductRepo struct{ mock.Mock }

func (m *MockProductRepo) Create(ctx context.Context, name string) (uint, error) {
	args := m.Called(ctx, name)
	return uint(args.Int(0)), args.Error(1)
}

func (m *MockProductRepo) List(ctx context.Context) ([]*entity.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *MockUow) WithTransaction(ctx context.Context, fn func(tx UnitOfWork) error) error {
	// In a unit test, we just call fn with ourselves
	return fn(m)
}
func (m *MockUow) Item() ItemRepo       { return m.Called().Get(0).(ItemRepo) }
func (m *MockUow) List() ListRepo       { return m.Called().Get(0).(ListRepo) }
func (m *MockUow) Product() ProductRepo { return m.Called().Get(0).(ProductRepo) }

type internalTestSuite struct {
	suite.Suite
	listRepo *MockListRepo
	itemRepo *MockItemRepo
	prodRepo *MockProductRepo
	uow      *MockUow
	listUc   IListUseCase
	itemUc   IItemUseCase
	ctx      context.Context
}

func (s *internalTestSuite) SetupSubTest() {
	s.ctx = context.Background()
	s.listRepo = new(MockListRepo)
	s.itemRepo = new(MockItemRepo)
	s.prodRepo = new(MockProductRepo)
	s.uow = new(MockUow)
	s.listUc = NewListUseCase(s.listRepo, s.itemRepo, s.uow)
	s.itemUc = NewItemUseCase(s.itemRepo, s.prodRepo, s.uow)
	s.uow.On("Item").Return(s.itemRepo)
	s.uow.On("List").Return(s.listRepo)
	s.uow.On("Product").Return(s.prodRepo)
}

func TestInternal(t *testing.T) {
	suite.Run(t, new(internalTestSuite))
}
