package headaches

import (
	"time"

	"encore.app/entity"
	"encore.dev/beta/errs"
)

func (repo *HeadacheRepository) ListHeadaches() entity.Headaches {
	var headaches entity.Headaches
	repo.db.Find(&headaches)
	return headaches
}

func (repo *HeadacheRepository) CreateHeadache(headache *entity.Headache) error {
	return repo.db.Create(headache).Error
}

func (repo *HeadacheRepository) DeleteHeadache(haId uint) error {
	tx := repo.db.Delete(&entity.Headache{}, haId)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return &errs.Error{Code: errs.NotFound}
	}
	return nil
}

func (repo *HeadacheRepository) GetHeadache(haId uint) (entity.Headache, error) {
	headache := entity.Headache{ID: haId}
	rows := repo.db.Find(&headache).RowsAffected
	if rows == 0 {
		return headache, &errs.Error{Code: errs.NotFound}
	}
	return headache, nil
}

func (repo *HeadacheRepository) PatchHeadache(
	haId uint,
	date *time.Time,
	severity *entity.HeadacheSeverity,
	types *entity.HeadacheTypes,
	positions *entity.HeadachePositions,
	symptoms *entity.HeadacheSymptoms,
) error {
	var headache = entity.Headache{ID: haId}
	if date != nil {
		headache.Date = *date
	}
	if severity != nil {
		headache.Severity = *severity
	}
	if types != nil {
		headache.Types = *types
	}
	if positions != nil {
		headache.Positions = *positions
	}
	if symptoms != nil {
		headache.Symptoms = *symptoms
	}
	tx := repo.db.Debug().Model(&headache).Updates(&headache)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return &errs.Error{Code: errs.NotFound}
	}
	return nil
}
