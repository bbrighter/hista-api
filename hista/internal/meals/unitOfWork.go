package meals

import (
	"context"

	"gorm.io/gorm"
)

type UnitOfWork struct {
	db *gorm.DB
	m  *MealRepository
	i  *ingredientRepo
	t  *templateRepo
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{
		db: db,
		m:  NewMealRepository(db),
		i:  newIngredientRepo(db),
		t:  newTemplateRepo(db),
	}
}

func (uow *UnitOfWork) Transaction(ctx context.Context, trans func(uow *UnitOfWork) error) error {
	return uow.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txUow := NewUnitOfWork(tx)
		return trans(txUow)
	})
}
