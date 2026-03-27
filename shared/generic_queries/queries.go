package generic_queries

import (
	"context"

	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type Piider interface {
	SetPiid(uuid.UUID)
}

func List[T Piider](ctx context.Context, db *gorm.DB) ([]T, error) {
	piid, err := PiidFromCtx(ctx)
	if err != nil {
		return []T{}, err
	}
	return gorm.G[T](db.Session(&gorm.Session{SkipDefaultTransaction: true})).Where("pi_id = ?", piid).Find(ctx)
}

func First[T Piider](ctx context.Context, db *gorm.DB, id uint) (T, error) {
	var t T
	piid, err := PiidFromCtx(ctx)
	if err != nil {
		return t, err
	}
	return gorm.G[T](db.Session(&gorm.Session{SkipDefaultTransaction: true})).Where("id = ?", id).Where("pi_id = ?", piid).First(ctx)
}

func Create[T Piider](ctx context.Context, db *gorm.DB, t T) error {
	piid, err := PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	t.SetPiid(piid)
	return gorm.G[T](db).Create(ctx, &t)
}

func Delete[T Piider](ctx context.Context, db *gorm.DB, id uint) error {
	return BatchDelete[T](ctx, db, []uint{id})
}

func BatchDelete[T Piider](ctx context.Context, db *gorm.DB, ids []uint) error {
	piid, err := PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	rows, err := gorm.G[T](db).Where("id IN ?", ids).Where("pi_id = ?", piid).Delete(ctx)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func Count[T Piider](ctx context.Context, db *gorm.DB, id uint) (int64, error) {
	piid, err := PiidFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	return gorm.G[T](db.Session(&gorm.Session{SkipDefaultTransaction: true})).Where("pi_id = ?", piid).Where("id = ?", id).Count(ctx, "*")
}

func UpdateColumn[T Piider](ctx context.Context, db *gorm.DB, id uint, columnName string, newValue any) error {
	piid, err := PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	rows, err := gorm.G[T](db).
		Where("id = ?", id).
		Where("pi_id = ?", piid).
		Update(ctx, columnName, newValue)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
