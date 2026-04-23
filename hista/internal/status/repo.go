package status

import (
	"context"

	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

type statusRepo struct {
	db *gorm.DB
}

func newStatusRepo(db *gorm.DB) *statusRepo {
	return &statusRepo{db: db}
}

func (r *statusRepo) ListStatus(ctx context.Context) ([]*Status, error) {
	return generic_queries.List[*Status](ctx, r.db)
}

func (r *statusRepo) CreateStatus(ctx context.Context, s *Status) error {
	return generic_queries.Create(ctx, r.db, s)
}

func (r *statusRepo) DeleteStatus(ctx context.Context, id uint) error {
	return generic_queries.Delete[*Status](ctx, r.db, id)
}

func (r *statusRepo) UpdateStatus(ctx context.Context, id uint, values map[string]any) error {
	return generic_queries.Updates(ctx, r.db, "statuses", id, values)
}
