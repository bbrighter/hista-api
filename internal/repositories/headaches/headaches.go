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
	tx := repo.db.Model(&entity.Headache{ID: haId}).Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return &errs.Error{Code: errs.NotFound}
	}
	return nil
}
