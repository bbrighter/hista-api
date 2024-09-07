package symptoms

import (
	"encore.app/entity"
	"encore.app/errors"
	"gorm.io/gorm/clause"
)

// Get all conditions including their conditionTypes
func (repo *SymptomsRepo) ListConditions(eventId uint) entity.Conditions {
	var conditions entity.Conditions
	repo.db.Where(&entity.Condition{ConditionEventID: eventId}).Preload(clause.Associations).Find(&conditions)
	return conditions
}

// Create a new condition based on the name. New symptoms are only created if the name in the corresponding category doesn't exist.
// Requires a conditionEventID
func (repo *SymptomsRepo) CreateConditionBySymptomName(condition *entity.Condition, symptomName string, categoryId uint) error {
	if condition.ConditionEventID == 0 {
		return errors.ErrorAttributeMustBeSet("ConditionEventID")
	}
	events := repo.db.Find(&entity.ConditionEvent{ID: condition.ConditionEventID}).RowsAffected
	if events == 0 {
		return errors.ErrorNotFound
	}
	var symptom = &entity.Symptom{Name: symptomName, SymptomCategoryID: categoryId}
	var err error

	symptomId, err := repo.CreateOrReplace(symptomName, categoryId)
	if err != nil {
		return err
	}
	symptom.ID = symptomId
	condition.Symptom = *symptom
	err = repo.db.Create(&condition).Error
	return err
}

// Create a new condition by SymptomID.
// Requires a ConditionEventID and SymptomID
func (repo *SymptomsRepo) CreateConditionBySymptomID(condition *entity.Condition) error {
	if condition.ConditionEventID == 0 || condition.SymptomID == 0 {
		return errors.ErrorAttributeMustBeSet("ConditionEventID and SymptomID")
	}
	events := repo.db.Find(&entity.ConditionEvent{ID: condition.ConditionEventID}).RowsAffected
	if events == 0 {
		return errors.ErrorNotFound
	}
	var symptom entity.Symptom
	symptoms := repo.db.Find(&symptom, &entity.Symptom{ID: condition.SymptomID}).RowsAffected
	if symptoms == 0 {
		return errors.ErrorNotFound
	}
	condition.Symptom = symptom
	var err error = repo.db.Create(condition).Error
	return err
}

// Delete a condition. Must contain ID.
// If the condition was the last one using a symptom, the symptom is deleted as well.
func (repo *SymptomsRepo) DeleteCondition(conditionId uint) error {
	var condition = entity.Condition{ID: conditionId}
	repo.db.Preload(clause.Associations).Find(&condition)
	tx := repo.db.Delete(&condition)
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}

func (repo *SymptomsRepo) ChangeSeverity(conditionId uint, newSeverity entity.ConditionSeverity) error {
	var condition = entity.Condition{ID: conditionId}
	tx := repo.db.Model(condition).Where(&condition).Updates(&entity.Condition{Severity: newSeverity})
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}

func (repo *SymptomsRepo) GetCondition(id uint) (entity.Condition, error) {
	var condition = entity.Condition{ID: id}
	err := repo.db.Preload(clause.Associations).First(&condition).Error
	return condition, err
}

// func numberToSeverity(no uint8) (entity.ConditionSeverity, error) {
// 	var err error
// 	var severity entity.ConditionSeverity
// 	switch no {
// 	case 1:
// 		severity = entity.VeryLow
// 	case 2:
// 		severity = entity.Low
// 	case 3:
// 		severity = entity.Medium
// 	case 4:
// 		severity = entity.High
// 	case 5:
// 		severity = entity.VeryHigh
// 	default:
// 		err = errors.NewError("invalid severity: "+string(no), 400)
// 	}
// 	return severity, err
// }
