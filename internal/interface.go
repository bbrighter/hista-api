package internal

import (
	"time"

	"encore.app/entity"
)

type (
	IMealsRepository interface {
		ListMeals() entity.Meals
		CreateMeal(date time.Time) (entity.Meal, error)
		GetMeal(id uint) (entity.Meal, error)
		DeleteMeal(entity.Meal) error
		PatchMeal(id uint, date *time.Time, freshness *entity.Freshness, stressLevel *uint8, isAlone *bool) error
	}

	IMealUseCase interface {
		ListMeals() entity.Meals
		CreateMeal(date *time.Time) (uint, error)
		GetMeal(id uint) (entity.Meal, error)
		DeleteMeal(id uint) error
		PatchMeal(id uint, date *time.Time, freshness *entity.Freshness, stressLevel *uint8, isAlone *bool) error
	}

	IFoodRepository interface {
		ListFoods(mealId uint) entity.Foods
		CreateFoodByName(mealId uint, ingredientName string) (uint, error)
		CreateFoodByID(mealId uint, ingredientId uint) (uint, error)
		DeleteFood(foodId uint) error
		ChangeCondition(foodId uint, condition entity.FoodCondition) error
		GetFood(foodId uint) entity.Food
	}

	IFoodUseCase interface {
		ListFoods(mealId uint) entity.Foods
		CreateFood(mealId uint, ingredientName string, ingredientId uint) (entity.Food, entity.Ingredients, error)
		DeleteFood(foodId uint) (entity.Ingredients, error)
		ChangeCondition(foodId uint, newCond entity.FoodCondition) error
	}

	IIngredientRepository interface {
		ListIngredients() entity.Ingredients
	}

	INotesRepository interface {
		List() entity.Notes
		Create() (entity.Note, error)
		Patch(note entity.Note, date *time.Time, text *string) error
		Delete(note entity.Note) error
	}

	INotesUseCase interface {
		List() entity.Notes
		Create() (uint, error)
		Patch(id uint, date *time.Time, text *string) error
		Delete(id uint) error
	}

	IConditionEventRepo interface {
		ListConditionEvents() entity.ConditionEvents
		CreateConditionEvent(date time.Time) (uint, error)
		GetConditionEvent(uint) (entity.ConditionEvent, error)
		PatchConditionEvent(uint, time.Time) error
		DeleteConditionEvent(uint) error
	}

	IConditionEventUseCase interface {
		List() entity.ConditionEvents
		Create(time.Time) (uint, error)
		Get(uint) (entity.ConditionEvent, error)
		Patch(uint, time.Time) error
		Delete(uint) error
	}

	IConditionRepo interface {
		ListConditions(eventId uint) entity.Conditions
		CreateConditionBySymptomName(eventId uint, symptomName string, symptomCategoryId uint) (uint, error)
		CreateConditionBySymptomID(eventId uint, symptomId uint) (uint, error)
		DeleteCondition(conditionId uint) error
		ChangeSeverity(conditionId uint, newSeverity entity.ConditionSeverity) error
		GetCondition(id uint) (entity.Condition, error)
	}

	IConditionUseCase interface {
		List(eventId uint) entity.Conditions
		Create(eventId uint, symptomName *string, symptomId *uint, symptomCategoryId *uint) (entity.Condition, entity.SymptomCategories, error)
		Delete(id uint) (entity.SymptomCategories, error)
		PatchSeverity(id uint, newSeverity entity.ConditionSeverity) error
	}

	ISymtpomsRepo interface {
		// ListSymptoms() entity.Symptoms
		CreateOrReplace(symptomName string, symtpomCategoryId uint) (uint, error)
	}

	ISymptomCategoriesRepo interface {
		ListCategories() entity.SymptomCategories
	}
)
