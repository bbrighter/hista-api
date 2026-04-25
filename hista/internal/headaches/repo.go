package headaches

import (
	"context"

	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

type HeadacheRepo struct {
	db *gorm.DB
}

func NewHeadacheRepo(db *gorm.DB) *HeadacheRepo {
	return &HeadacheRepo{db: db}
}

func (r *HeadacheRepo) ListHeadaches(ctx context.Context) ([]*Headache, error) {
	return generic_queries.List[*Headache](ctx, r.db)
}

func (r *HeadacheRepo) FirstHeadache(ctx context.Context, id uint) (*Headache, error) {
	return generic_queries.First[*Headache](ctx, r.db, id)
}

func (r *HeadacheRepo) DeleteHeadache(ctx context.Context, id uint) error {
	return generic_queries.Delete[*Headache](ctx, r.db, id)
}

func (r *HeadacheRepo) CreateHeadache(ctx context.Context, headache *Headache) error {
	return generic_queries.Create(ctx, r.db, headache)
}

func (r *HeadacheRepo) UpdateHeadache(ctx context.Context, id uint, values map[string]any) error {
	return generic_queries.Updates(ctx, r.db, "headaches", id, values)
}
