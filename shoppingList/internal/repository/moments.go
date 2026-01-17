package repository

import (
	"context"
	"errors"
	"time"

	"encore.app/shared/generic_queries"
	"encore.app/shoppingList/entity"
	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type MomentRepo struct {
	db *gorm.DB
	moments map[uuid.UUID]time.Time
}

func NewMomentRepo(db *gorm.DB) MomentRepo {
	return MomentRepo{db: db, moments: make(map[uuid.UUID]time.Time)}
}

func (r MomentRepo) GetOrCreate(ctx context.Context) (*entity.Moment, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return &entity.Moment{}, err
	}
	if moment, ok := r.getCache(piid); ok {
		return moment, nil
	}
	moment, err := gorm.G[*entity.Moment](r.db).Where("pi_id = ?", piid).First(ctx)
	if err == nil {
		r.updateCache(moment)
		return moment, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		moment = &entity.Moment{PIID: piid}
		err := gorm.G[*entity.Moment](r.db).Create(ctx, &moment)
		r.updateCache(moment)
		return moment, err
	}

	r.updateCache(moment)
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

	r.moments[piid] = moment.UpdatedAt
	return err
}


func (r MomentRepo) updateCache(moment *entity.Moment) {
	r.moments[moment.PIID] = moment.UpdatedAt
}

func (r MomentRepo) getCache(piid uuid.UUID) (*entity.Moment, bool) {
	if updatedAt, ok := r.moments[piid]; ok {
		return &entity.Moment{PIID: piid, UpdatedAt: updatedAt}, true
	}
	return nil, false
}