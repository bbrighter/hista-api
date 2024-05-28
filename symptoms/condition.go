package symptoms

import (
	"encore.app/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Condition struct {
	ID               uint
	Symptom          Symptom
	SymptomID        uint
	Severity         ConditionSeverity
	ConditionEventID uint
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

// Get all conditions including their conditionTypes
func getConditions(service *Service, eventID uint) Conditions {
	var conditions Conditions
	service.db.Where(&Condition{ConditionEventID: eventID}).Preload(clause.Associations).Find(&conditions)
	return conditions
}

// Create a new condition based on the name. New symptoms are only created if the name in the corresponding category doesn't exist.
// Requires a conditionEventID
func (condition *Condition) createConditionBySymptomName(service *Service, symptomName string, symptomCategoryId uint) (SymptomCategories, error) {
	if condition == nil {
		return SymptomCategories{}, errors.ErrorNil
	}
	if condition.ConditionEventID == 0 {
		return SymptomCategories{}, errors.ErrorAttributeMustBeSet("ConditionEventId")
	}
	var symptom = &Symptom{Name: symptomName, SymptomCategoryID: symptomCategoryId}
	var err error
	err = symptom.createOrReplace(service)
	if err != nil {
		return SymptomCategories{}, err
	}
	condition.Symptom = *symptom
	condition.Severity = Medium

	err = service.db.Create(condition).Error
	var symptoms SymptomCategories = getSymptomCategories(service)
	return symptoms, err
}

// Create a new condition by SymptomID.
// Requires a ConditionEventID and SymptomID
func (condition *Condition) createConditionBySymptomID(service *Service) error {
	if condition == nil {
		return errors.ErrorNil
	}
	if condition.SymptomID == 0 {
		return errors.ErrorAttributeMustBeSet("SymptomID")
	}
	if condition.ConditionEventID == 0 {
		return errors.ErrorAttributeMustBeSet("ConditionEventID")
	}
	condition.Severity = Medium
	var err error = service.db.Create(condition).Error
	service.db.Preload(clause.Associations).Find(condition)
	return err
}

// Delete a condition. Must contain ID.
// If the condition was the last one using a symptom, the symptom is deleted as well.
func (condition *Condition) delete(service *Service) error {
	if condition == nil {
		return errors.ErrorNil
	}
	if condition.ID == 0 {
		return errors.ErrorAttributeMustBeSet("id")
	}
	service.db.Preload(clause.Associations).Find(condition)
	var err error = service.db.Transaction(func(tx *gorm.DB) error {
		var symptom Symptom = condition.Symptom
		if err := tx.Delete(condition).Error; err != nil {
			return err
		}
		return symptom.deleteIfUnused(tx)
	})
	return err
}

func (condition *Condition) changeSeverity(service *Service, newSeverity ConditionSeverity) error {
	if condition.ID == 0 {
		return errors.ErrorAttributeMustBeSet("id")
	}
	condition.Severity = newSeverity
	tx := service.db.Model(condition).Updates(Condition{Severity: newSeverity})
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
