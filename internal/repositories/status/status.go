package status

import (
	"time"

	"encore.app/entity"
	"encore.app/errors"
)

func (repo *StatusRepo) Create(status *entity.Status) error {
	return repo.db.Create(status).Error
}

func (repo *StatusRepo) First(status *entity.Status) error {
	if rows := repo.db.First(&status).RowsAffected; rows == 0 {
		return errors.ErrorNotFound
	}
	return nil
}

func (repo *StatusRepo) Update(status *entity.Status) error {
	return repo.db.Model(&entity.Status{}).Where(&entity.Status{ID: status.ID}).Updates(status).Error
}

func (repo *StatusRepo) Find() (statuses entity.Statuses) {
	repo.db.Find(&statuses)
	return statuses
}

func (repo *StatusRepo) FindForDate(date time.Time) (entity.Status, bool) {
	var status entity.Status
	tx := repo.db.Where("DATE(date) = ?", date.Format("2006-01-02")).First(&status)
	return status, tx.RowsAffected > 0
}

func (repo *StatusRepo) Delete(status *entity.Status) error {
	tx := repo.db.Delete(status)
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}
