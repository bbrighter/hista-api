package status

import (
	"time"

	"encore.app/entity"
	"encore.app/errors"
	"gorm.io/gorm/clause"
)

func (repo *StatusRepo) Save(status *entity.Status) error {
	err := repo.db.Save(status).Error
	return err
}

func (repo *StatusRepo) Find() (statuses entity.Statuses) {
	repo.db.Find(&statuses)
	return statuses
}

func (repo *StatusRepo) FindForDate(date time.Time) (entity.Status, bool) {
	var status entity.Status
	tx := repo.db.Where("DATE(date) = ?", date.Format("2006-01-02")).Preload(clause.Associations).First(&status)
	return status, tx.RowsAffected > 0
}

func (repo *StatusRepo) Delete(status *entity.Status) error {
	tx := repo.db.Delete(status)
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}
