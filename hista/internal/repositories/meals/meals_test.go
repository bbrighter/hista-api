package meals

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"encore.dev/types/uuid"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const GUID_STR = "cf0d4408-8db5-4572-b5d9-4ed873d1341f"

var GUID = uuid.FromStringOrNil(GUID_STR)

func initTest(t *testing.T) (*MealRepository, context.Context) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(
		&entity.Meal{},
		&entity.Ingredient{},
		&entity.Food{},
	)
	assert.NoError(t, err)

	ctx := context.WithValue(t.Context(), "piid", GUID)
	return &MealRepository{db: db}, ctx
}

func TestCreateMeal(t *testing.T) {
	repo, ctx := initTest(t)

	var err error

	var meal = new(entity.Meal)
	meal.Date = time.Date(1999, 0, 0, 0, 0, 0, 0, time.Local)
	err = repo.CreateMeal(ctx, meal)

	assert.GreaterOrEqual(t, meal.ID, uint(1))
	assert.NoError(t, err)
	var mealInDB entity.Meal
	repo.db.First(&mealInDB, entity.Meal{ID: meal.ID})
	assert.True(t, mealInDB.Date.Equal(mealInDB.Date))
}

func TestGetMeals(t *testing.T) {
	repo, ctx := initTest(t)

	meals, err := repo.ListMeals(ctx)
	assert.NoError(t, err)
	assert.Equal(t, len(meals), 0)
}

func TestGetMeal(t *testing.T) {
	repo, ctx := initTest(t)

	var err error

	_, err = repo.GetMeal(ctx, 1000)
	assert.Error(t, err)

	var meal = entity.Meal{Freshness: entity.Fresh}
	err = generic_queries.Create(ctx, repo.db, &meal)
	require.NoError(t, err)
	var ing = entity.Ingredient{Name: "name"}
	err = generic_queries.Create(ctx, repo.db, &ing)
	require.NoError(t, err)
	var food = entity.Food{MealID: meal.ID, IngredientID: ing.ID}
	err = generic_queries.Create(ctx, repo.db, &food)
	require.NoError(t, err)

	meal, err = repo.GetMeal(ctx, meal.ID)
	assert.NoError(t, err)
	assert.Equal(t, entity.Fresh, meal.Freshness)
	assert.Equal(t, len(meal.Foods), 1)
}

func TestDeleteMeal(t *testing.T) {
	repo, ctx := initTest(t)

	var err error
	err = repo.DeleteMeal(ctx, 100)
	assert.Error(t, err)

	var meal = entity.Meal{Freshness: entity.Fresh}
	err = generic_queries.Create(ctx, repo.db, &meal)
	require.NoError(t, err)
	var ing = entity.Ingredient{PIID: GUID, Name: "name"}
	err = generic_queries.Create(ctx, repo.db, &ing)
	require.NoError(t, err)
	var food = entity.Food{PIID: GUID, MealID: meal.ID, MealPIID: GUID, IngredientID: ing.ID, IngredientPIID: GUID}
	err = generic_queries.Create(ctx, repo.db, &food)
	require.NoError(t, err)

	err = repo.DeleteMeal(ctx, meal.ID)
	assert.NoError(t, err)

	foods := repo.db.Find(&entity.Food{}).RowsAffected
	assert.EqualValues(t, 0, foods)
	ingredients := repo.db.Find(&entity.Ingredient{}).RowsAffected
	assert.EqualValues(t, 0, ingredients)
}

func TestPatchMeal(t *testing.T) {
	repo, ctx := initTest(t)

	var err error

	var meal = entity.Meal{Freshness: entity.Fresh, PIID: GUID, StressLevel: 1}
	err = gorm.G[entity.Meal](repo.db).Create(ctx, &meal)
	require.NoError(t, err)
	var ing = entity.Ingredient{PIID: GUID, Name: "name"}
	err = gorm.G[entity.Ingredient](repo.db).Create(ctx, &ing)
	require.NoError(t, err)
	var food = entity.Food{PIID: GUID, MealID: meal.ID, MealPIID: GUID, IngredientID: ing.ID, IngredientPIID: GUID}
	err = gorm.G[entity.Food](repo.db).Create(ctx, &food)
	require.NoError(t, err)

	var patchDate time.Time = time.Date(1700, 0, 0, 0, 0, 0, 0, time.Local)
	var patchFreshness entity.Freshness = entity.Older

	err = repo.PatchMeal(ctx, meal.ID, &patchDate, &patchFreshness, nil, nil)
	assert.NoError(t, err)

	mealInDb, err := gorm.G[entity.Meal](repo.db).Where("id = ?", meal.ID).First(ctx)
	assert.NoError(t, err)
	assert.True(t, mealInDb.Date.Equal(patchDate))
	assert.Equal(t, mealInDb.Freshness, patchFreshness)
	assert.EqualValues(t, mealInDb.StressLevel, 1)

	err = repo.PatchMeal(ctx, 1000, nil, nil, nil, nil)
	assert.Error(t, err)

	// var meal Meal
	// var err error
	// err = meal.create(repo)
	// defer meal.delete(repo)
	// assert.NoError(t, err)

	// var params = PatchParams{}
	// var patchDate time.Time = time.Date(1700, 0, 0, 0, 0, 0, 0, time.Local)
	// var patchFreshness Freshness = Older
	// params.Date = &patchDate
	// err = meal.patch(repo, params)
	// assert.NoError(t, err)
	// params.Freshness = &patchFreshness
	// err = meal.patch(repo, params)
	// assert.NoError(t, err)
	// var patchStressLevel uint8
	// params.StressLevel = &patchStressLevel
	// err = meal.patch(repo, params)
	// assert.NoError(t, err)

	// var mealInDB = Meal{ID: meal.ID}
	// repo.db.First(&mealInDB)
	// assert.True(t, mealInDB.Date.Equal(patchDate))
	// assert.Equal(t, patchFreshness, mealInDB.Freshness)
	// assert.Equal(t, patchStressLevel, mealInDB.StressLevel)
}

// func TestGetMealsAndDependencies(t *testing.T) {
// 	repo, _ := initTest(t)

// 	var err error
// 	_, err = GetMealsAndDependencies(repo.db)
// 	assert.NoError(t, err)

// 	repo.db.Create(&entity.Meal{
// 		ID:          100,
// 		Freshness:   entity.Fresh,
// 		StressLevel: 1,
// 		Foods: []entity.Food{
// 			{ID: 1,
// 				Ingredient: entity.Ingredient{ID: 10, Name: "ingredient"}},
// 		},
// 	})

// 	var meals entity.Meals
// 	meals, err = GetMealsAndDependencies(repo.db)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "ingredient", meals[0].Foods[0].Ingredient.Name)
// }
