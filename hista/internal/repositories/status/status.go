package status

import (
	"context"
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
)

func (repo *StatusRepo) Create(ctx context.Context, status *entity.Status) error {
	return generic_queries.Create(ctx, repo.db, status)
}

func (repo *StatusRepo) First(ctx context.Context, id uint) (*entity.Status, error) {
	return generic_queries.First[*entity.Status](ctx, repo.db, id)
}

func (repo *StatusRepo) Update(ctx context.Context, status *entity.Status) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	return repo.db.Model(&entity.Status{}).Where(&entity.Status{ID: status.ID}).Where("pi_id = ?", piid).Updates(status).Error
}

func (repo *StatusRepo) Find(ctx context.Context) ([]*entity.Status, error) {
	return generic_queries.List[*entity.Status](ctx, repo.db)
}

func (repo *StatusRepo) FindForDate(ctx context.Context, date time.Time) (entity.Status, bool) {
	var status entity.Status
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return status, false
	}
	tx := repo.db.Where("pi_id = ?", piid).Where("DATE(date) = ?", date.Format("2006-01-02")).First(&status)
	return status, tx.RowsAffected > 0
}

func (repo *StatusRepo) Delete(ctx context.Context, id uint) error {
	return generic_queries.Delete[*entity.Status](ctx, repo.db, id)
}
