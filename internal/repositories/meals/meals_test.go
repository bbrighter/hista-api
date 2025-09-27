package meals

import (
	"context"
	"testing"
	"time"

	"encore.app/entity"
	"encore.dev/types/uuid"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const GUID_STR = "cf0d4408-8db5-4572-b5d9-4ed873d1341f"

func initTest(t *testing.T) (*MealRepository, context.Context) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(
		&entity.Meal{},
		&entity.Ingredient{},
		&entity.Food{},
	)
	assert.NoError(t, err)

	ctx := context.WithValue(t.Context(), "piid", uuid.FromStringOrNil(GUID_STR))
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
	var meal entity.Meal

	meal, err = repo.GetMeal(ctx, 1000)
	assert.Error(t, err)

	repo.db.Create(&entity.Meal{
		ID:        100,
		Freshness: entity.Fresh,
		PIID:      uuid.FromStringOrNil(GUID_STR),
		Foods: []entity.Food{
			{ID: 1},
		},
	})

	meal, err = repo.GetMeal(ctx, 100)
	assert.NoError(t, err)
	assert.Equal(t, entity.Fresh, meal.Freshness)
	assert.Equal(t, len(meal.Foods), 1)
}

func TestDeleteMeal(t *testing.T) {
	repo, ctx := initTest(t)

	var err error
	var mealId uint = 100
	err = repo.DeleteMeal(ctx, mealId)
	assert.Error(t, err)

	repo.db.Create(&entity.Meal{
		ID:        mealId,
		Freshness: entity.Fresh,
		PIID:      uuid.FromStringOrNil(GUID_STR),
		Foods: []entity.Food{
			{ID: 1, PIID: uuid.FromStringOrNil(GUID_STR),
				Ingredient: entity.Ingredient{ID: 1, PIID: uuid.FromStringOrNil(GUID_STR)}},
		},
	})

	err = repo.DeleteMeal(ctx, mealId)
	assert.NoError(t, err)

	foods := repo.db.Find(&entity.Food{}).RowsAffected
	assert.EqualValues(t, 0, foods)
	ingredients := repo.db.Find(&entity.Ingredient{}).RowsAffected
	assert.EqualValues(t, 0, ingredients)
}

func TestPatchMeal(t *testing.T) {
	repo, ctx := initTest(t)

	repo.db.Create(&entity.Meal{
		ID:          100,
		Freshness:   entity.Fresh,
		StressLevel: 1,
		PIID:        uuid.FromStringOrNil(GUID_STR),
		Foods: []entity.Food{
			{ID: 1},
		},
	})

	var patchDate time.Time = time.Date(1700, 0, 0, 0, 0, 0, 0, time.Local)
	var patchFreshness entity.Freshness = entity.Older

	var err error
	err = repo.PatchMeal(ctx, 100, &patchDate, &patchFreshness, nil, nil)
	assert.NoError(t, err)

	var mealInDb entity.Meal
	repo.db.Find(&mealInDb, &entity.Meal{ID: 100})
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
