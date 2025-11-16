package repository

import (
	"context"
	"errors"
	"time"

	"encore.app/shared/generic_queries"
	"encore.app/shoppingList/entity"
	"gorm.io/gorm"
)

type MomentRepo struct {
	db *gorm.DB
}

func NewMomentRepo(db *gorm.DB) MomentRepo {
	return MomentRepo{db: db}
}

func (r MomentRepo) GetOrCreate(ctx context.Context) (*entity.Moment, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return &entity.Moment{}, err
	}
	moment, err := gorm.G[*entity.Moment](r.db).Where("pi_id = ?", piid).First(ctx)
	if err == nil {
		return moment, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		moment = &entity.Moment{PIID: piid}
		err := gorm.G[*entity.Moment](r.db).Create(ctx, &moment)
		return moment, err
	}

	return moment, err
}

func (r MomentRepo) Update(ctx context.Context) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	moment := &entity.Moment{PIID: piid, UpdatedAt: time.Now()}

	rows, err := gorm.G[entity.Moment](r.db).Where("pi_id = ?", piid).Update(ctx, "updated_at", moment.UpdatedAt)
	if err == nil && rows == 0 {
		return gorm.ErrRecordNotFound
	}

	return err

}
