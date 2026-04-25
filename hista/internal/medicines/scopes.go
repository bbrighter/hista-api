package medicines

import (
	"context"

	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

func isNotOld(tx *gorm.DB) *gorm.DB {
	return tx.Where("intakes.date > current_date - 7")
}

func gIsNotOld(tx *gorm.Statement) {
	tx.Where("intakes.date > current_date - 7")
}

func wherePiid(ctx context.Context) func(tx *gorm.Statement) {
	return func(tx *gorm.Statement) {
		piid, err := generic_queries.PiidFromCtx(ctx)
		if err != nil {
			tx.AddError(err)
			return
		}
		tx.Where("pi_id = ?", piid)
	}
}
