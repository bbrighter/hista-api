package generic_queries

import (
	"context"

	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type Piider interface {
	SetPiid(uuid.UUID)
}

func WherePiid(ctx context.Context) func(tx *gorm.Statement) {
	return func(tx *gorm.Statement) {
		piid, err := PiidFromCtx(ctx)
		if err != nil {
			tx.AddError(err)
			return
		}
		tx.Where("pi_id = ?", piid)
	}
}

func List[T Piider](ctx context.Context, db *gorm.DB) ([]T, error) {
	return gorm.G[T](db.Session(&gorm.Session{SkipDefaultTransaction: true})).
		Scopes(WherePiid(ctx)).
		Find(ctx)
}

func First[T Piider](ctx context.Context, db *gorm.DB, id uint) (T, error) {
	return gorm.G[T](db.Session(&gorm.Session{SkipDefaultTransaction: true})).
		Where("id = ?", id).
		Scopes(WherePiid(ctx)).
		First(ctx)
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
	rows, err := gorm.G[T](db).
		Where("id IN ?", ids).
		Scopes(WherePiid(ctx)).
		Delete(ctx)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func Count[T Piider](ctx context.Context, db *gorm.DB, id uint) (int64, error) {
	return gorm.G[T](db.Session(&gorm.Session{SkipDefaultTransaction: true})).
		Scopes(WherePiid(ctx)).
		Where("id = ?", id).
		Count(ctx, "*")
}

func UpdateColumn[T Piider](ctx context.Context, db *gorm.DB, id uint, columnName string, newValue any) error {
	rows, err := gorm.G[T](db).
		Where("id = ?", id).
		Scopes(WherePiid(ctx)).
		Update(ctx, columnName, newValue)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func Updates(ctx context.Context, db *gorm.DB, tableName string, id uint, values map[string]any) error {
	rows, err := gorm.G[map[string]any](db).
		Table(tableName).
		Where("id = ?", id).
		Scopes(WherePiid(ctx)).
		Updates(ctx, values)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
