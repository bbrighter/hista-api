package meals

import (
	"time"

	"encore.app/errors"
	"encore.dev/beta/errs"
	"gorm.io/gorm"
)

type Freshness uint8

const (
	Fresh   Freshness = 0
	SameDay Freshness = 1
	Older   Freshness = 2
)

type Meal struct {
	ID          uint
	Date        time.Time
	Freshness   Freshness
	StressLevel uint8
	IsAlone     bool
	Foods       []Food `gorm:"constraint:OnDelete:CASCADE"`
}

type Meals []Meal

func (meals *Meals) get(service *Service) error {
	return service.db.Find(meals).Error
}

func (meal *Meal) create(service *Service) error {
	if meal.Date.IsZero() {
		meal.Date = time.Now()
	}
	return service.db.Create(&meal).Error
}

func (meal *Meal) get(service *Service) error {
	if meal.ID == 0 {
		return errors.ErrorIDMissing
	}
	if service.db.Preload("Foods.Ingredient").Preload("Foods").Find(&meal).RowsAffected == 0 {
		return &errs.Error{Code: errs.NotFound}
	}
	return nil
}

func (meal *Meal) delete(service *Service) error {
	if meal.ID == 0 {
		return errors.ErrorIDMissing
	}
	if service.db.Find(&Meal{ID: meal.ID}).RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return service.db.Transaction(func(tx *gorm.DB) error {
		var foods []Food
		tx.Where(&Food{MealID: meal.ID}).Find(&foods)
		if err := tx.Delete(&Food{}, Food{MealID: meal.ID}).Error; err != nil {
			return err
		}
		for _, food := range foods {
			deleteIngredientIfUnused(tx, food.IngredientID)
		}
		return tx.Delete(&Meal{ID: meal.ID}).Error
	})
}

type PatchParams struct {
	Date        *time.Time
	Freshness   *Freshness
	StressLevel *uint8
	IsAlone     *bool
}

// Patch a meal with parameters. Only given parameters are patched.
func (meal *Meal) patch(service *Service, params PatchParams) error {
	if meal.ID == 0 {
		return errors.ErrorIDMissing
	}
	if rows := service.db.First(&meal, &Meal{ID: meal.ID}).RowsAffected; rows == 0 {
		return errors.ErrorNotFound
	}
	tx := service.db.Model(&Meal{ID: meal.ID})
	var updates = make(map[string]interface{})
	if params.Date != nil {
		meal.Date = *params.Date
		updates["date"] = *params.Date
	}
	if params.Freshness != nil {
		meal.Freshness = *params.Freshness
		updates["freshness"] = *params.Freshness
	}
	if params.StressLevel != nil {
		meal.StressLevel = *params.StressLevel
		updates["stress_level"] = *params.StressLevel
	}
	if params.IsAlone != nil {
		meal.IsAlone = *params.IsAlone
		updates["is_alone"] = *params.IsAlone
	}
	return tx.Updates(updates).Error
}

func GetMealsAndDependencies(db *gorm.DB) (Meals, error) {
	var meals Meals
	var err error = db.Preload("Foods.Ingredient").
		Preload("Foods").
		Find(&meals).Error
	return meals, err
}
