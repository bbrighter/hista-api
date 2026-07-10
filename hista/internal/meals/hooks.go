package meals

import "gorm.io/gorm"

func (f *Food) AfterDelete(tx *gorm.DB) error {
	if f.IngredientID == 0 {
		return nil
	}
	var count int64
	if err := tx.Model(&Food{}).
		Where("ingredient_id = ?", f.IngredientID).
		Count(&count).Error; err != nil {
		return err
	}
	print("COUNT:", count)
	if count > 0 {
		return nil
	}
	return tx.Delete(&Ingredient{}, f.IngredientID).Error
}
