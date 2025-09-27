package meals

import (
	"gorm.io/gorm"
)

type MealRepository struct {
	db *gorm.DB
}

func NewMealRepository(db *gorm.DB) *MealRepository {
	return &MealRepository{db: db}
}

// func piidFromCtx(ctx context.Context) (uuid.UUID, error) {
// 	piid, ok := ctx.Value("piid").(uuid.UUID)
// 	if !ok {
// 		return uuid.Nil, errors.New("missing PIID in context")
// 	}
// 	return piid, nil
// }

// func scopedDB[T any](ctx context.Context, db *gorm.DB) (gorm.ChainInterface[T], error) {
// 	piid, err := piidFromCtx(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return gorm.G[T](db).Where("pi_id = ?", piid), nil
// }
