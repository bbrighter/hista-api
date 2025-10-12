package headaches

import (
	"context"
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"encore.dev/beta/errs"
)

func (repo *HeadacheRepository) ListHeadaches(ctx context.Context) ([]*entity.Headache, error) {
	return generic_queries.List[*entity.Headache](ctx, repo.db)
}

func (repo *HeadacheRepository) CreateHeadache(ctx context.Context, headache *entity.Headache) error {
	return generic_queries.Create(ctx, repo.db, headache)
}

func (repo *HeadacheRepository) DeleteHeadache(ctx context.Context, haId uint) error {
	return generic_queries.Delete[*entity.Headache](ctx, repo.db, haId)
}

func (repo *HeadacheRepository) GetHeadache(ctx context.Context, haId uint) (*entity.Headache, error) {
	return generic_queries.First[*entity.Headache](ctx, repo.db, haId)
}

func (repo *HeadacheRepository) PatchHeadache(
	ctx context.Context,
	haId uint,
	date *time.Time,
	severity *entity.HeadacheSeverity,
	types *entity.HeadacheTypes,
	positions *entity.HeadachePositions,
	symptoms *entity.HeadacheSymptoms,
	description *string,
) error {
	updates := make(map[string]interface{})
	if date != nil {
		updates["date"] = *date
	}
	if severity != nil {
		updates["severity"] = *severity
	}
	if types != nil {
		updates["types"] = *types
	}
	if positions != nil {
		updates["positions"] = *positions
	}
	if symptoms != nil {
		updates["symptoms"] = *symptoms
	}
	if description != nil {
		updates["description"] = *description
	}
	if len(updates) == 0 {
		return nil
	}
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	tx := repo.db.Model(&entity.Headache{ID: haId}).Where("pi_id = ?", piid).Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return &errs.Error{Code: errs.NotFound}
	}
	return nil
}
