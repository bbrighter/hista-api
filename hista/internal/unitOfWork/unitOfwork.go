package unitofwork

import (
	"context"

	"encore.app/hista/internal"
	"encore.app/hista/internal/repositories/medicine"
	"gorm.io/gorm"
)

type unitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) internal.UnitOfWork {
	return unitOfWork{db: db}
}

func (u unitOfWork) WithTransaction(ctx context.Context, fn func(tx internal.UnitOfWork) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		r := NewUnitOfWork(tx)
		return fn(r)
	})
}

func (u unitOfWork) Medicine() internal.IMedicineRepo {
	return medicine.NewMedicineRepo(u.db)
}
