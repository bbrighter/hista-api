package status

import (
	"time"

	"encore.app/entity"
	"encore.app/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *StatusRepo) Create(status *entity.Status) error {
	return repo.db.Create(status).Error
}

func (repo *StatusRepo) Update(status *entity.Status) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		if !status.Date.IsZero() {
			if err := tx.Model(status).Update("date", status.Date).Error; err != nil {
				return err
			}
		}
		if status.Evening != nil {
			if err := tx.Save(status.Evening).Error; err != nil {
				return err
			}
		}
		if status.Morning != nil {
			if err := tx.Save(status.Morning).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (repo *StatusRepo) Save(status *entity.Status) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(status).Error; err != nil {
			return err
		}
		if status.Morning != nil {
			if err := tx.Save(status.Morning).Error; err != nil {
				return err
			}
		}
		if status.Evening != nil {
			if err := tx.Save(status.Evening).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (repo *StatusRepo) Find() (statuses entity.Statuses) {
	repo.db.Preload(clause.Associations).Find(&statuses)
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
