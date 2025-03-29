package internal

import (
	"testing"
	"time"

	"encore.app/entity"
)

type (
	IMealsRepository interface {
		ListMeals() entity.Meals
		CreateMeal(*entity.Meal) error
		GetMeal(id uint) (entity.Meal, error)
		DeleteMeal(id uint) error
		PatchMeal(id uint, date *time.Time, freshness *entity.Freshness, stressLevel *uint8, isAlone *bool) error
		ListMealsWithDependencies() entity.Meals
	}

	IMealUseCase interface {
		List() entity.Meals
		Create(date *time.Time) (entity.Meal, error)
		Get(id uint) (entity.Meal, error)
		Delete(id uint) (entity.Ingredients, error)
		Patch(id uint, date *time.Time, freshness *entity.Freshness, stressLevel *uint8, isAlone *bool) error
	}

	IFoodRepository interface {
		ListFoods(mealId uint) entity.Foods
		CreateFoodByName(food *entity.Food, ingredientName string) error
		CreateFoodByID(food *entity.Food) error
		DeleteFood(foodId uint) error
		ChangeCondition(foodId uint, condition entity.FoodCondition) error
		GetFood(foodId uint) entity.Food
	}

	IFoodUseCase interface {
		List(mealId uint) entity.Foods
		Create(mealId uint, ingredientName string, ingredientId uint) (entity.Food, entity.Ingredients, error)
		Delete(foodId uint) (entity.Ingredients, error)
		ChangeCondition(foodId uint, newCond entity.FoodCondition) error
	}

	IIngredientRepository interface {
		ListIngredients() entity.Ingredients
	}

	IIngredientUseCase interface {
		List() entity.Ingredients
	}

	INotesRepository interface {
		List() entity.Notes
		Create(note *entity.Note) error
		Patch(id uint, date *time.Time, text *string) error
		Delete(id uint) error
	}

	INotesUseCase interface {
		List() entity.Notes
		Create() (entity.Note, error)
		Patch(id uint, date *time.Time, text *string) error
		Delete(id uint) error
	}

	IConditionEventRepo interface {
		ListConditionEvents() entity.ConditionEvents
		CreateConditionEvent(event *entity.ConditionEvent) error
		GetConditionEvent(uint) (entity.ConditionEvent, error)
		PatchConditionEvent(uint, time.Time) error
		DeleteConditionEvent(uint) error
		ListConditionEventsAndDependencies() entity.ConditionEvents
	}

	IConditionEventUseCase interface {
		List() entity.ConditionEvents
		Create() (entity.ConditionEvent, error)
		Get(uint) (entity.ConditionEvent, error)
		Patch(uint, time.Time) error
		Delete(uint) (entity.SymptomCategories, error)
	}

	IConditionRepo interface {
		ListConditions(eventId uint) entity.Conditions
		CreateConditionBySymptomName(condition *entity.Condition, symptomName string, symptomCategoryId uint) error
		CreateConditionBySymptomID(condition *entity.Condition) error
		DeleteCondition(conditionId uint) error
		ChangeSeverity(conditionId uint, newSeverity entity.Severity) error
		GetCondition(id uint) (entity.Condition, error)
	}

	IConditionUseCase interface {
		List(eventId uint) entity.Conditions
		Create(eventId uint, symptomName *string, symptomId *uint, symptomCategoryId *uint) (entity.Condition, entity.SymptomCategories, error)
		Delete(id uint) (entity.SymptomCategories, error)
		PatchSeverity(id uint, newSeverity entity.Severity) error
	}

	ISymptomsRepo interface {
		CreateOrReplace(symptomName string, symptomCategoryId uint) (uint, error)
	}

	ISymptomCategoriesRepo interface {
		ListCategories() entity.SymptomCategories
		CreateCategory(*entity.SymptomCategory) error
	}

	ISymptomsUseCase interface {
		PutSymptom(symptomName string, symptomCategoryId uint) (uint, error)
		List() entity.SymptomCategories
		CreateCategory(name string) (uint, error)
	}

	IDiaryUseCase interface {
		Get() (entity.Meals, entity.ConditionEvents, entity.SymptomCategories, entity.Notes, entity.PollenEvents)
	}

	IStatisticsRepo interface {
		FindSymptomsForFoods(fromDate time.Time, toDate time.Time, ingredientIds []uint) (entity.FoodResults, error)
		FindFoodForSymptoms(fromDate time.Time, toDate time.Time, symptomIds []uint) (entity.SymptomResults, error)
		CountFoods(symptomIds []uint) []entity.CountResult
		CountSymptoms(ingredientIds []uint) []entity.CountResult
	}

	IStatisticsUseCase interface {
		FindSymptomsForFoods(fromDate time.Time, toDate time.Time, ingredientIds []uint) (entity.FoodResults, error)
		FindFoodForSymptoms(fromDate time.Time, toDate time.Time, symptomIds []uint) (entity.SymptomResults, error)
	}

	IPollenRepo interface {
		FindPollenWithSeverity(severity int) entity.PollenEvents
		Create(pollen entity.Pollens, lastUpdated time.Time) error
	}

	IDWDRepo interface {
		GetKarlsruheData() (entity.DWDPollen, error)
		DwdStringToDate() (time.Time, error)
		UseTestQuery(*testing.T)
	}

	IPollenUseCase interface {
		List() entity.PollenEvents
		Create() error
		UseTestQuery(*testing.T)
	}

	IHeadacheRepo interface {
		ListHeadaches() []entity.Headache
		CreateHeadache(ha *entity.Headache) error
		DeleteHeadache(haId uint) error
		GetHeadache(haId uint) (entity.Headache, error)
		PatchHeadache(haId uint, date *time.Time, severity *entity.Severity, types *entity.HeadacheTypes, positions *entity.HeadachePositions, symptoms *entity.HeadacheSymptoms)
	}
)
