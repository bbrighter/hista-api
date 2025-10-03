package entity

// func (e *ConditionEvent) BeforeDelete(tx *gorm.DB) (err error) {
// 	var conditions Conditions
// 	cons := tx.Preload(clause.Associations).
// 		Where("pi_id = ?", e.PIID).
// 		Where("condition_event_id = ?", e.ID).
// 		Find(&conditions)
// 	if cons.RowsAffected == 0 {
// 		return nil
// 	}
// 	return tx.Where("pi_id = ?", e.PIID).Delete(&conditions).Error
// }

// func (c *Condition) AfterDelete(tx *gorm.DB) (err error) {
// 	return deleteSymptomIfUnused(tx, c.SymptomID, c.Symptom.SymptomCategoryID, c.PIID)
// }

// func deleteCategoryIfUnused(tx *gorm.DB, catId uint, piid uuid.UUID) error {
// 	var err error
// 	result := tx.Debug().Where("pi_id = ?", piid).
// 		Where("symptom_category_id = ?", catId).
// 		Find(&Symptom{})
// 	rows := result.RowsAffected
// 	err = result.Error
// 	if rows == 0 {
// 		err = tx.Where("pi_id = ?", piid).Delete(&SymptomCategory{ID: catId}).Error
// 	}
// 	return err
// }

// func deleteSymptomIfUnused(tx *gorm.DB, symptomId uint, catId uint, piid uuid.UUID) error {
// 	var err error
// 	var rows int64
// 	if symptomId == 0 || catId == 0 {
// 		return nil
// 	}
// 	usedConditions := tx.Preload(clause.Associations).Find(&Condition{}, &Condition{SymptomID: symptomId, PIID: piid})
// 	if usedConditions.RowsAffected == 0 {
// 		result := tx.Where("pi_id = ?", piid).Where(&Symptom{ID: symptomId}).Delete(&Symptom{})
// 		err = result.Error
// 		rows = result.RowsAffected
// 	}
// 	if rows > 0 && err == nil {
// 		err = deleteCategoryIfUnused(tx, catId, piid)
// 	}
// 	return err
// }

// func (m *Meal) BeforeDelete(tx *gorm.DB) error {
// 	var foods Foods
// 	tx.Where("pi_id = ?", m.PIID).Where("meal_id = ?", m.ID).Find(&foods)
// 	for _, food := range foods {
// 		err := tx.Delete(&food).Error
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (f *Food) AfterDelete(tx *gorm.DB) (err error) {
// 	var unusedIngredients Ingredients
// 	resp := tx.Table("ingredients").
// 		Joins("LEFT JOIN foods ON foods.ingredient_id = ingredients.id").
// 		Where("foods.id IS NULL").
// 		Where("pi_id = ?", f.PIID).
// 		Find(&unusedIngredients)
// 	if resp.RowsAffected > 0 {
// 		err = tx.Where("pi_id = ?", f.PIID).Delete(&unusedIngredients).Error
// 	}
// 	return err
// }

// func (pollenEvent *PollenEvent) BeforeDelete(tx *gorm.DB) error {
// 	return tx.Delete(&Pollen{}, &Pollen{PollenEventID: pollenEvent.ID}).Error
// }
