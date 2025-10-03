package move

import (
	"context"

	"encore.app/entity"
	"encore.app/generic_queries"
	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type MoveRepo struct {
	db *gorm.DB
}

func NewMoveRepo(db *gorm.DB) MoveRepo {
	return MoveRepo{db: db}
}

func (r MoveRepo) Move(ctx context.Context, fromPiid, toPiid uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := movePiid[*entity.Headache](ctx, tx, fromPiid, toPiid); err != nil {
			return err
		}
		if err := movePiid[*entity.Meal](ctx, tx, fromPiid, toPiid); err != nil {
			return err
		}
		if err := movePiid[*entity.Food](ctx, tx, fromPiid, toPiid); err != nil {
			return err
		}
		if err := movePiid[*entity.Ingredient](ctx, tx, fromPiid, toPiid); err != nil {
			return err
		}
		if err := movePiid[*entity.ConditionEvent](ctx, tx, fromPiid, toPiid); err != nil {
			return err
		}
		if err := movePiid[*entity.Condition](ctx, tx, fromPiid, toPiid); err != nil {
			return err
		}
		if err := movePiid[*entity.Symptom](ctx, tx, fromPiid, toPiid); err != nil {
			return err
		}
		if err := movePiid[*entity.SymptomCategory](ctx, tx, fromPiid, toPiid); err != nil {
			return err
		}
		if err := movePiid[*entity.Note](ctx, tx, fromPiid, toPiid); err != nil {
			return err
		}
		if err := movePiid[*entity.Status](ctx, tx, fromPiid, toPiid); err != nil {
			return err
		}
		return nil
	})
}

func movePiid[T generic_queries.Piider](ctx context.Context, tx *gorm.DB, fromPiid, toPiid uuid.UUID) error {
	var err error
	if fromPiid == uuid.Nil {
		gorm.G[T](tx).Where("pi_id IS NULL").Update(ctx, "pi_id", toPiid)
	} else {
		gorm.G[T](tx).Where("pi_id = ?", fromPiid).Update(ctx, "pi_id", toPiid)
	}
	return err
}
