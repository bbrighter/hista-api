package symptoms

import (
	"context"

	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

func wherePiid(ctx context.Context) func(db *gorm.Statement) {
	return func(db *gorm.Statement) {
		piid, err := generic_queries.PiidFromCtx(ctx)
		if err != nil {
			db.AddError(err)
		}
		db.Where("pi_id = ?", piid)
	}
}
