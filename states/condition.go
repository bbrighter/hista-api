package states

import (
	"encore.app/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Condition struct {
	ID        uint
	Symptom   Symptom
	SymptomID uint
	Severity  ConditionSeverity
	StateID   uint
}

type Conditions []Condition

type ConditionSeverity uint8

const (
	VeryLow  ConditionSeverity = 1
	Low      ConditionSeverity = 2
	Medium   ConditionSeverity = 3
	High     ConditionSeverity = 4
	VeryHigh ConditionSeverity = 5
)

func newCondition(severity ConditionSeverity, stateID uint) *Condition {
	return &Condition{
		Severity: severity,
		StateID:  stateID,
	}
}

// Get all conditions including their conditionTypes
func getConditions(service *Service, stateID uint) Conditions {
	var conditions Conditions
	service.db.Where(&Condition{StateID: stateID}).Preload(clause.Associations).Find(&conditions)
	return conditions
}

// Create a new condition.
// Requires a severity and stateId
func (condition *Condition) create(service *Service, symptomName string, symptomCategoryId uint) error {
	if condition == nil {
		return errors.ErrorNil
	}
	if condition.Severity == 0 || condition.Severity > 5 {
		return errors.ErrorAttributeMustBeSet("severity")
	}
	if condition.StateID == 0 {
		return errors.ErrorAttributeMustBeSet("stateId")
	}
	var symptom = &Symptom{Name: symptomName, SymptomCategoryID: symptomCategoryId}
	var err error
	err = symptom.createOrReplace(service)
	if err != nil {
		return err
	}
	condition.Symptom = *symptom

	return service.db.Create(condition).Error
}

func (condition *Condition) delete(service *Service) error {
	if condition == nil {
		return errors.ErrorNil
	}
	if condition.ID == 0 {
		return errors.ErrorAttributeMustBeSet("id")
	}
	service.db.Preload("ConditionType").Find(condition)
	var err error = service.db.Transaction(func(tx *gorm.DB) error {
		var conditionType Symptom = condition.Symptom
		if err := tx.Delete(condition).Error; err != nil {
			return err
		}
		return conditionType.deleteIfUnused(tx)
	})
	return err
}

func (condition *Condition) changeSeverity(service *Service, newSeverity ConditionSeverity) error {
	if condition.ID == 0 {
		return errors.ErrorAttributeMustBeSet("id")
	}
	condition.Severity = newSeverity
	tx := service.db.Debug().Model(condition).Updates(Condition{Severity: newSeverity})
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}

func numberToSeverity(no uint8) (ConditionSeverity, error) {
	var err error
	var severity ConditionSeverity
	switch no {
	case 1:
		severity = VeryLow
	case 2:
		severity = Low
	case 3:
		severity = Medium
	case 4:
		severity = High
	case 5:
		severity = VeryHigh
	default:
		err = errors.NewError("invalid severity: "+string(no), 400)
	}
	return severity, err
}
