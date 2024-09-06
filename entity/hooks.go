package entity

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (e *ConditionEvent) AfterDelete(tx *gorm.DB) (err error) {
	var conditions Conditions
	cons := tx.Preload(clause.Associations).Find(&conditions, Condition{ConditionEventID: e.ID})
	if cons.RowsAffected == 0 {
		return nil
	}
	return tx.Delete(&conditions).Error
}

func (c *Condition) AfterDelete(tx *gorm.DB) (err error) {
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

func (m *Meal) BeforeDelete(tx *gorm.DB) error {
	var foods Foods
	tx.Find(&foods, Food{MealID: m.ID})
	for _, food := range foods {
		err := tx.Delete(&food).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Food) AfterDelete(tx *gorm.DB) error {
	var unusedIngredients Ingredients
	tx.Table("ingredients").
		Joins("LEFT JOIN foods ON foods.ingredient_id = ingredients.id").
		Where("foods.id IS NULL").
		Find(&unusedIngredients)
	return tx.Delete(&unusedIngredients).Error
}

func (pollenEvent *PollenEvent) BeforeDelete(tx *gorm.DB) error {
	return tx.Delete(&Pollen{}, &Pollen{PollenEventID: pollenEvent.ID}).Error
}
