package pollen

import (
	"context"
	"time"

	"encore.app/errors"
	"gorm.io/gorm"
)

type PollenRepo struct {
	db *gorm.DB
}

func NewPollenRepo(db *gorm.DB) *PollenRepo {
	return &PollenRepo{db: db}
}

func (r *PollenRepo) CreatePollen(ctx context.Context, p *PollenEvent) error {
	return gorm.G[PollenEvent](r.db).Create(ctx, p)
}

func (r *PollenRepo) ListPollen(ctx context.Context) ([]PollenEvent, error) {
	return gorm.G[PollenEvent](r.db).Preload("Pollens", nil).Find(ctx)
}

func (r *PollenRepo) DoesExistAfter(ctx context.Context, time time.Time) error {
	count, err := gorm.G[PollenEvent](r.db).Where("created_at > ?", time).Count(ctx, "*")
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.ErrorAlreadyExists
	}
	return nil
}

func (r *PollenRepo) ListPollenWithSeverity(ctx context.Context, severity int) ([]PollenEvent, error) {
	return gorm.G[PollenEvent](r.db).
		Preload("Pollens", func(db gorm.PreloadBuilder) error {
			db.Where("intensity > ?", severity)
			return nil
		}).Find(ctx)
}
