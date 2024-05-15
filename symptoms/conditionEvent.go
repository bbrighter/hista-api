package symptoms

import (
	"time"

	"encore.app/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConditionEvent struct {
	ID         uint
	Date       time.Time
	Conditions []Condition
}

type ConditionEvents []ConditionEvent

func newConditionEvent(date time.Time) *ConditionEvent {
	return &ConditionEvent{Date: date}
}

func (event *ConditionEvent) create(service *Service) error {
	if event == nil {
		return errors.ErrorNil
	}
	if event.Date.IsZero() {
		return errors.ErrorAttributeMustBeSet("date")
	}
	return service.db.Create(event).Error
}

func getConditionEvents(service *Service) ConditionEvents {
	var events ConditionEvents
	service.db.Find(&events)
	return events
}

func getConditionEvent(service *Service, id uint) (ConditionEvent, error) {
	var event = ConditionEvent{ID: id}
	if service.db.Preload("Conditions.Symptom").Preload(clause.Associations).Find(&event).RowsAffected == 0 {
		return event, errors.ErrorNotFound
	}
	return event, nil
}

func (event *ConditionEvent) delete(service *Service) error {
	if event.ID == 0 {
		return errors.ErrorIDMissing
	}
	if service.db.Find(event).RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	var err error = service.db.Transaction(func(tx *gorm.DB) error {
		var conditions Conditions
		tx.Where(&Condition{ConditionEventID: event.ID}).Preload(clause.Associations).Find(&conditions)
		if err := tx.Delete(&Condition{}, Condition{ConditionEventID: event.ID}).Error; err != nil {
			return err
		}
		for _, condition := range conditions {
			if err := condition.Symptom.deleteIfUnused(tx); err != nil {
				return err
			}
		}
		return tx.Delete(event).Error
	})
	return err
}

func (event *ConditionEvent) patch(service *Service, date time.Time) error {
	if event.ID == 0 {
		return errors.ErrorIDMissing
	}
	if service.db.Find(event).RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return service.db.Model(event).Update("date", date).Error
}

func GetConditionEventsAndDependencies(db *gorm.DB) ConditionEvents {
	var events ConditionEvents
	db.Preload("Conditions.Symptom").Preload("Conditions").Find(&events)
	return events
}
