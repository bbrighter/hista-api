package entity

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (e *ConditionEvent) AfterDelete(tx *gorm.DB) (err error) {
	var conditions Conditions
	cons := tx.Debug().Preload(clause.Associations).Find(&conditions, Condition{ConditionEventID: e.ID})
	if cons.RowsAffected == 0 {
		return nil
	}
	return tx.Debug().Delete(&conditions).Error
}

func (c *Condition) AfterDelete(tx *gorm.DB) (err error) {
	log.Printf("Condition: ID: %v, Severity: %v, SymptomId: %v, ConditionEventId: %v, Symptom.ID: %v, Symptom.Name: %v, Symtpom.CategoryId: %v",
		c.ID,
		c.Severity,
		c.SymptomID,
		c.ConditionEventID,
		c.Symptom.ID,
		c.Symptom.Name,
		c.Symptom.SymptomCategoryID)
	return deleteSymptomIfUnused(tx, c.SymptomID, c.Symptom.SymptomCategoryID)
}

func deleteCategoryIfUnused(tx *gorm.DB, catId uint) error {
	var err error
	result := tx.Find(&Symptom{}, &Symptom{SymptomCategoryID: catId})
	rows := result.RowsAffected
	err = result.Error
	if rows == 0 {
		err = tx.Delete(&SymptomCategory{ID: catId}).Error
	}
	return err
}

func deleteSymptomIfUnused(tx *gorm.DB, symptomId uint, catId uint) error {
	var err error
	var rows int64
	if symptomId == 0 || catId == 0 {
		return nil
	}
	usedConditions := tx.Preload(clause.Associations).Find(&Condition{}, &Condition{SymptomID: symptomId})
	if usedConditions.RowsAffected == 0 {
		result := tx.Where(&Symptom{ID: symptomId}).Delete(&Symptom{})
		err = result.Error
		rows = result.RowsAffected
	}
	if rows > 0 && err == nil {
		err = deleteCategoryIfUnused(tx, catId)
	}
	return err
}
