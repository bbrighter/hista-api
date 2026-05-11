package moments

import (
	"context"
	"sync"

	"encore.app/shared/generic_queries"
	"encore.dev/types/uuid"
)

type MomentRepo struct {
	mu       sync.RWMutex
	items    map[uuid.UUID]int
	products map[uuid.UUID]int
}

func NewMomentRepo() *MomentRepo {
	return &MomentRepo{items: make(map[uuid.UUID]int), products: make(map[uuid.UUID]int)}
}

func (r *MomentRepo) UpdateItems(ctx context.Context) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	r.mu.Lock()
	r.items[piid] += 1
	r.mu.Unlock()
	return nil
}

func (r *MomentRepo) UpdateProducts(ctx context.Context) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	r.mu.Lock()
	r.products[piid] += 1
	r.mu.Unlock()
	return nil
}

func (r *MomentRepo) GetMoments(ctx context.Context) (Moment, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return Moment{}, err
	}
	moment := Moment{PIID: piid}

	r.mu.RLock()
	defer r.mu.RUnlock()
	version, ok := r.items[piid]
	if ok {
		moment.ItemsVersion = version
	}

	version, ok = r.products[piid]
	if ok {
		moment.ProductsVersion = version
	}

	return moment, nil
}
