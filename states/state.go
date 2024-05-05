package states

import (
	"time"

	"encore.app/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type State struct {
	ID         uint
	Date       time.Time
	Conditions []Condition
}

type States []State

func newState(date time.Time) *State {
	return &State{Date: date}
}

func (state *State) create(service *Service) error {
	if state == nil {
		return errors.ErrorNil
	}
	if state.Date.IsZero() {
		return errors.ErrorAttributeMustBeSet("date")
	}
	return service.db.Create(state).Error
}

func getStates(service *Service) States {
	var states States
	service.db.Find(&states)
	return states
}

func getState(service *Service, id uint) (State, error) {
	var state = State{ID: id}
	if service.db.Preload("Conditions.ConditionType").Preload("Conditions").Find(&state).RowsAffected == 0 {
		return state, errors.ErrorNotFound
	}
	return state, nil
}

func (state *State) delete(service *Service) error {
	if state.ID == 0 {
		return errors.ErrorIDMissing
	}
	if service.db.Find(state).RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	var err error = service.db.Transaction(func(tx *gorm.DB) error {
		var conditions Conditions
		tx.Where(&Condition{StateID: state.ID}).Preload(clause.Associations).Find(&conditions)
		if err := tx.Delete(&Condition{}, Condition{StateID: state.ID}).Error; err != nil {
			return err
		}
		for _, condition := range conditions {
			if err := condition.Symptom.deleteIfUnused(tx); err != nil {
				return err
			}
		}
		return tx.Delete(state).Error
	})
	return err
}
