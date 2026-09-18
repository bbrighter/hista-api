package meals

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

const DefaultFoodCondition FoodCondition = Cooked

type MealService struct {
	m   *MealRepository
	t   *templateRepo
	i   *ingredientRepo
	uow *UnitOfWork
}

func NewMealService(db *gorm.DB) *MealService {
	m := NewMealRepository(db)
	t := newTemplateRepo(db)
	i := newIngredientRepo(db)
	uow := NewUnitOfWork(db)
	return &MealService{m: m, t: t, i: i, uow: uow}
}

func (uc MealService) ListMeals(ctx context.Context) (Meals, error) {
	return uc.m.ListAllMeals(ctx)
}
func (uc MealService) CreateMeal(ctx context.Context, date time.Time) (Meal, error) {
	var meal = &Meal{
		Date:        date,
		IsAlone:     true,
		Freshness:   Fresh,
		StressLevel: 0,
	}
	_, err := uc.m.CreateMeal(ctx, meal)
	return *meal, err
}
func (uc MealService) GetMeal(ctx context.Context, id uint) (Meal, error) {
	meal, err := uc.m.GetMealAndFoods(ctx, id)
	if err != nil {
		return Meal{}, err
	}
	return *meal, nil
}
func (uc MealService) DeleteMeal(ctx context.Context, id uint) (Ingredients, error) {
	if err := uc.m.DeleteMeal(ctx, id); err != nil {
		return Ingredients{}, err
	}
	return uc.i.ListIngredients(ctx)
}
func (uc MealService) UpdateMeal(ctx context.Context, id uint, date *time.Time, freshness *uint8, stressLevel *uint8, isAlone *bool) error {
	values := make(map[string]any)
	if date != nil {
		values["date"] = *date
	}
	if freshness != nil {
		values["freshness"] = Freshness(*freshness)
	}
	if stressLevel != nil {
		values["stress_level"] = *stressLevel
	}
	if isAlone != nil {
		values["is_alone"] = *isAlone
	}
	return uc.m.UpdateMeal(ctx, id, values)
}

func (uc MealService) ListFood(ctx context.Context, mealId uint) (Foods, error) {
	return uc.m.ListFoodsByMeal(ctx, mealId)
}
func (uc MealService) CreateFoodByName(ctx context.Context, mealId uint, ingredientName string) (Food, error) {
	ingredientName = strings.TrimSpace(ingredientName)
	var food = Food{
		Condition: DefaultFoodCondition,
		MealID:    mealId,
	}
	if ingredientName == "" {
		return Food{}, errors.New("ingredient name cannot be empty")
	}
	ingredient, err := uc.i.GetIngredientByName(ctx, ingredientName)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err := uc.m.CreateFoodAndIngredient(ctx, &food, ingredientName)
		return food, err
	} else if err != nil {
		return Food{}, err
	} else {
		food.IngredientID = ingredient.ID
		uc.m.CreateFood(ctx, &food)
		return food, err
	}
}

func (uc MealService) CreateFoodById(ctx context.Context, mealId uint, ingredientId uint) (Food, error) {
	food := Food{MealID: mealId, IngredientID: ingredientId, Condition: DefaultFoodCondition}
	err := uc.uow.Transaction(ctx, func(uow *UnitOfWork) error {
		if err := uow.m.CreateFood(ctx, &food); err != nil {
			return err
		}
		if err := uow.i.UpdateIngredient(ctx, ingredientId, map[string]any{"is_archived": false}); err != nil {
			return err
		}

		return nil
	})
	return food, err
}

func (uc MealService) DeleteFood(ctx context.Context, foodId uint) error {
	return uc.m.DeleteFood(ctx, foodId)
}
func (uc MealService) ChangeFoodCondition(ctx context.Context, foodId uint, newCond string) error {
	values := map[string]any{"condition": newCond}
	return uc.m.UpdateFood(ctx, foodId, values)
}
func (uc MealService) ChangeFoodAmount(ctx context.Context, foodId uint, amount *int) error {
	values := map[string]any{"amount": amount}
	return uc.m.UpdateFood(ctx, foodId, values)
}

func (uc MealService) ListIngredients(ctx context.Context) (Ingredients, error) {
	return uc.i.ListIngredients(ctx)
}

func (uc MealService) UpdateIngredient(ctx context.Context, ingredientId uint, newName *string, newNutrition *Nutrition, isArchived *bool) error {
	values := make(map[string]any)
	if newName != nil {
		values["name"] = strings.TrimSpace(*newName)
	}
	if newNutrition != nil {
		values["nutrition_protein"] = newNutrition.Protein
		values["nutrition_carbohydrate"] = newNutrition.Carbohydrate
		values["nutrition_fat"] = newNutrition.Fat
		values["nutrition_fiber"] = newNutrition.Fiber
	}
	if isArchived != nil {
		values["is_archived"] = *isArchived
	}
	return uc.i.UpdateIngredient(ctx, ingredientId, values)
}

func (uc MealService) ToggleIngredientArchived(ctx context.Context, ingredientId uint) error {
	values := map[string]any{"is_archived": gorm.Expr("NOT is_archived")}
	return uc.i.UpdateIngredient(ctx, ingredientId, values)
}

func (uc MealService) ApplyTemplate(ctx context.Context, mealId uint, templateId uint) (Foods, error) {
	template, err := uc.t.FirstTemplateAndItems(ctx, templateId)
	if err != nil {
		return Foods{}, err
	}
	var foods = []Food{}
	for _, it := range template.Items {
		var food = Food{
			Condition:    it.Condition,
			IngredientID: it.IngredientID,
			MealID:       mealId,
		}
		foods = append(foods, food)
	}

	err = uc.m.CreateFoods(ctx, foods)

	var ids = []uint{}
	for i := range foods {
		ids = append(ids, foods[i].ID)
	}
	foods, err = uc.m.ListFoodsByIds(ctx, ids)
	return foods, err
}

func (uc MealService) ListTemplatesAndItems(ctx context.Context) (Templates, error) {
	return uc.t.ListTemplatesAndItems(ctx)
}

func (uc MealService) CreateTemplate(ctx context.Context, name string, items []TemplateItem) (uint, error) {
	template := Template{
		Name:  name,
		Items: items,
	}
	if err := uc.t.CreateTemplate(ctx, &template); err != nil {
		return 0, err
	}
	return template.ID, nil
}

func (uc MealService) ReplaceTemplate(ctx context.Context, id uint, name string, items []TemplateItem) error {
	return uc.t.ReplaceTemplate(ctx, id, name, items)
}

func (uc MealService) DeleteTemplate(ctx context.Context, id uint) error {
	return uc.t.DeleteTemplate(ctx, id)
}
