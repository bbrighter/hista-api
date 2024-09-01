package meals

import (
	"time"

	"encore.app/entity"
	"encore.app/errors"
	"encore.dev/beta/errs"
	"gorm.io/gorm"
)

const mealIsAlone = true

func (repo *MealRepository) ListMeals() entity.Meals {
	var meals entity.Meals
	repo.db.Find(&meals)
	return meals
}

func (repo *MealRepository) CreateMeal(data time.Time) (uint, error) {
	var meal = entity.Meal{IsAlone: mealIsAlone}
	if meal.Date.IsZero() {
		meal.Date = time.Now()
	}
	err := repo.db.Create(&meal).Error
	return meal.ID, err
}

func (repo *MealRepository) GetMeal(id uint) (entity.Meal, error) {
	var meal = entity.Meal{ID: id}
	if repo.db.Debug().Preload("Foods.Ingredient").Preload("Foods").Find(&meal).RowsAffected == 0 {
		return meal, &errs.Error{Code: errs.NotFound}
	}
	return meal, nil
}

func (repo *MealRepository) DeleteMeal(id uint) error {
	var meal = entity.Meal{ID: id}
	if repo.db.Find(&meal).RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return repo.db.Transaction(func(tx *gorm.DB) error {
		var foods []entity.Food
		tx.Where(&entity.Food{MealID: meal.ID}).Find(&foods)
		if err := tx.Delete(&entity.Food{}, entity.Food{MealID: meal.ID}).Error; err != nil {
			return err
		}
		for _, food := range foods {
			deleteIngredientIfUnused(tx, food.IngredientID)
		}
		return tx.Delete(&entity.Meal{ID: meal.ID}).Error
	})
}

// PatchMeal a meal with parameters. Only given parameters are patched.
func (repo *MealRepository) PatchMeal(
	id uint,
	date *time.Time,
	freshness *entity.Freshness,
	stressLevel *uint8,
	isAlone *bool,
) error {
	var meal entity.Meal
	if rows := repo.db.First(&meal, &entity.Meal{ID: id}).RowsAffected; rows == 0 {
		return errors.ErrorNotFound
	}
	tx := repo.db.Model(&entity.Meal{ID: meal.ID})
	var updates = make(map[string]interface{})
	if date != nil {
		meal.Date = *date
		updates["date"] = *date
	}
	if freshness != nil {
		meal.Freshness = *freshness
		updates["freshness"] = *freshness
	}
	if stressLevel != nil {
		meal.StressLevel = *stressLevel
		updates["stress_level"] = *stressLevel
	}
	if isAlone != nil {
		meal.IsAlone = *isAlone
		updates["is_alone"] = *isAlone
	}
	return tx.Updates(updates).Error
}

func GetMealsAndDependencies(db *gorm.DB) (entity.Meals, error) {
	var meals entity.Meals
	var err error = db.Preload("Foods.Ingredient").
		Preload("Foods").
		Find(&meals).Error
	return meals, err
}

func deleteIngredientIfUnused(tx *gorm.DB, ingredientId uint) error {
	var foods []entity.Food
	var err error
	if usedIngredients := tx.Where(entity.Food{IngredientID: ingredientId}).Find(&foods).RowsAffected; usedIngredients == 0 {
		err = tx.Where(&entity.Ingredient{ID: ingredientId}).Delete(&entity.Ingredient{}).Error
	}
	return err
}
