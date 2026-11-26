package meta

import (
	"context"
	"database/sql"
	"time"

	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

type MetaRepo struct {
	db *gorm.DB
}

func NewMetaRepo(db *gorm.DB) *MetaRepo {
	return &MetaRepo{db: db}
}

type CountEntities struct {
	Entity string `gorm:"column:entity"`
	Filled bool   `gorm:"column:filled"`
}

func (r *MetaRepo) AreTodaysEntitiesFilled(ctx context.Context) ([]CountEntities, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return []CountEntities{}, err
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	entities, err := gorm.G[CountEntities](r.db).Raw(`
		SELECT 'headaches' AS entity, count(*) > 0 AS filled 
		FROM headaches 
		WHERE pi_id = @piid AND date > @date

		UNION ALL

		SELECT 'statuses', count(*) > 0 
		FROM statuses 
		WHERE pi_id = @piid AND date > @date
		`,
		sql.Named("piid", piid),
		sql.Named("date", today),
	).Find(ctx)

	return entities, err
}
