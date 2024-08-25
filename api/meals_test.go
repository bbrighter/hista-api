package api

import (
	"context"
	"testing"
	"time"

	entity "encore.app/entity"
	"github.com/stretchr/testify/assert"
)

var testMeal *entity.Meal
var testFood *entity.Food
var testIngredient *entity.Ingredient

func (service *Service) initData() {
	var meal = &entity.Meal{
		Date:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
		Freshness:   entity.Fresh,
		StressLevel: 3,
		IsAlone:     true,
	}
	service.DB.FirstOrCreate(&meal, &meal)
	var ingredient = &entity.Ingredient{Name: "Ingredient"}
	service.DB.FirstOrCreate(&ingredient, &ingredient)
	var food = &entity.Food{IngredientID: ingredient.ID, Condition: entity.Cooked, MealID: meal.ID}
	service.DB.FirstOrCreate(&food, &food)

	testMeal = meal
	testFood = food
	testIngredient = ingredient
}

func initTest(t *testing.T) *Service {
	service, err := initService()
	assert.NoError(t, err)
	service.initData()
	return service
}

func initAPITest(t *testing.T) (*Service, context.Context) {
	var ctx context.Context = context.TODO()
	service := initTest(t)
	return service, ctx
}

func TestGetMealsAPI(t *testing.T) {
	service, ctx := initAPITest(t)
	service.initData()

	resp, err := service.GetMeals(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(resp.Meals), 1)
}

func TestPostMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var now = time.Now()
	var params = entity.MealParams{Date: &now}
	resp, err := service.PostMeal(ctx, params)
	defer service.DeleteMeal(ctx, resp.ID)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.ID, uint(1))
}

func TestGetMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	_, err = service.GetMeal(ctx, 100)
	assert.Error(t, err)

	var resp entity.MealResponse
	resp, err = service.GetMeal(ctx, testMeal.ID)
	assert.NoError(t, err, err)
	assert.EqualValues(t, testMeal.ID, resp.ID)
}

func TestDeleteMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	err = service.DeleteMeal(ctx, 10000)
	assert.Error(t, err)

	err = service.DeleteMeal(ctx, testMeal.ID)
	assert.NoError(t, err)
}

func TestPatchMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var meal = entity.Meal{ID: testMeal.ID}
	defer service.DeleteMeal(ctx, testMeal.ID)

	var params entity.MealParams
	var now time.Time = time.Now()
	params.Date = &now
	var err error
	err = service.PatchMeal(ctx, meal.ID, params)
	assert.NoError(t, err)

	var stressLevel uint8 = 2
	params.StressLevel = &stressLevel
	err = service.PatchMeal(ctx, meal.ID, params)
	assert.NoError(t, err)

	err = service.PatchMeal(ctx, 10000, params)
	assert.Error(t, err)
}

func TestGetFoods(t *testing.T) {
	service, ctx := initAPITest(t)
	service.initData()

	resp, err := service.GetFoods(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, resp.Foods, 0)
}

func TestPostFood(t *testing.T) {
	service, ctx := initAPITest(t)
	service.initData()

	var params = FoodParams{IngredientName: "New", Condition: entity.Cooked}
	_, err := service.PostFood(ctx, 1, params)
	assert.Error(t, err)
}

func TestDeleteFoodAPI(t *testing.T) {
	service, ctx := initAPITest(t)
	service.initData()

	var err error
	ing, err := service.DeleteFood(ctx, testMeal.ID, testFood.ID)
	assert.NoError(t, err)
	assert.Len(t, ing.Ingredients, 0)
}

func TestPatchFoodCondition(t *testing.T) {
	service, ctx := initAPITest(t)
	service.initData()

	var params = FoodConditionParams{Condition: "raw"}

	var err error
	service.PatchFoodCondition(ctx, testMeal.ID, testFood.ID, params)
	assert.NoError(t, err)

	var invalidParams = FoodConditionParams{Condition: "not valid"}
	service.PatchFoodCondition(ctx, testMeal.ID, testFood.ID, invalidParams)
	assert.NoError(t, err)
}
