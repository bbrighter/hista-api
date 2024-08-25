package meals

import (
	"testing"
	"time"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initTest(t *testing.T) *MealRepository {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&entity.Meal{}, &entity.Ingredient{}, &entity.Food{})
	assert.NoError(t, err)
	return &MealRepository{db: db}
}

func TestCreateMeal(t *testing.T) {
	repo := initTest(t)

	var err error

	time := time.Date(1999, 0, 0, 0, 0, 0, 0, time.Local)
	id, err := repo.Create(time)

	assert.GreaterOrEqual(t, id, uint(1))
	assert.NoError(t, err)
	var mealInDB entity.Meal
	repo.db.First(&mealInDB, entity.Meal{ID: id})
	assert.True(t, mealInDB.Date.Equal(mealInDB.Date))
}

func TestGetMeals(t *testing.T) {
	repo := initTest(t)

	meals := repo.List()
	assert.Equal(t, len(meals), 0)
}

func TestGetMeal(t *testing.T) {
	repo := initTest(t)

	var err error
	var meal entity.Meal

	meal, err = repo.Get(1000)
	assert.Error(t, err)

	repo.db.Create(&entity.Meal{
		ID:        100,
		Freshness: entity.Fresh,
		Foods: []entity.Food{
			{ID: 1},
		},
	})

	meal, err = repo.Get(100)
	assert.NoError(t, err)
	assert.Equal(t, entity.Fresh, meal.Freshness)
	assert.Equal(t, len(meal.Foods), 1)
}

func TestDeleteMeal(t *testing.T) {
	repo := initTest(t)

	var err error
	err = repo.Delete(100)
	assert.Error(t, err)

	repo.db.Create(&entity.Meal{
		ID:        100,
		Freshness: entity.Fresh,
		Foods: []entity.Food{
			{ID: 1},
		},
	})

	err = repo.Delete(100)
	assert.NoError(t, err)
}

func TestPatchMeal(t *testing.T) {
	repo := initTest(t)

	repo.db.Create(&entity.Meal{
		ID:          100,
		Freshness:   entity.Fresh,
		StressLevel: 1,
		Foods: []entity.Food{
			{ID: 1},
		},
	})

	var params entity.PatchParams
	var patchDate time.Time = time.Date(1700, 0, 0, 0, 0, 0, 0, time.Local)
	var patchFreshness entity.Freshness = entity.Older
	params.Date = &patchDate
	params.Freshness = &patchFreshness

	var err error
	err = repo.Patch(100, params)
	assert.NoError(t, err)

	var mealInDb entity.Meal
	repo.db.Find(&mealInDb, &entity.Meal{ID: 100})
	assert.True(t, mealInDb.Date.Equal(patchDate))
	assert.Equal(t, mealInDb.Freshness, patchFreshness)
	assert.EqualValues(t, mealInDb.StressLevel, 1)

	err = repo.Patch(1000, params)
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

func TestGetMealsAndDependencies(t *testing.T) {
	repo := initTest(t)

	var err error
	_, err = GetMealsAndDependencies(repo.db)
	assert.NoError(t, err)

	repo.db.Create(&entity.Meal{
		ID:          100,
		Freshness:   entity.Fresh,
		StressLevel: 1,
		Foods: []entity.Food{
			{ID: 1,
				Ingredient: entity.Ingredient{ID: 10, Name: "ingredient"}},
		},
	})

	var meals entity.Meals
	meals, err = GetMealsAndDependencies(repo.db)
	assert.NoError(t, err)
	assert.Equal(t, "ingredient", meals[0].Foods[0].Ingredient.Name)
}
