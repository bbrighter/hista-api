package meals

import (
	"testing"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetFood(t *testing.T) {
	repo := initTest(t)

	var foods entity.Foods
	foods = repo.ListFoods(10)
	assert.Len(t, foods, 0)

	repo.db.Create(&entity.Food{ID: 1, MealID: 10})

	foods = repo.ListFoods(10)
	assert.Len(t, foods, 1)
}

func TestCreateFoodByName(t *testing.T) {
	repo := initTest(t)

	var err error
	// err = repo.CreateFoodByName(1, "New name")
	// assert.Error(t, err)

	repo.db.Create(&entity.Meal{ID: 1})
	err = repo.CreateFoodByName(1, "New name")
	assert.NoError(t, err)
}

func TestCreateFoodById(t *testing.T) {
	repo := initTest(t)

	var err error
	err = repo.CreateFoodByID(1, 10)
	assert.Error(t, err)

	repo.db.Create(&entity.Meal{ID: 1})
	repo.db.Create(&entity.Ingredient{ID: 2})
	err = repo.CreateFoodByID(1, 2)
	assert.NoError(t, err)
}

func TestDeleteFood(t *testing.T) {
	repo := initTest(t)

	var err error
	err = repo.DeleteFood(1)
	assert.Error(t, err)

	repo.db.Create(&entity.Food{ID: 1, Ingredient: entity.Ingredient{ID: 10}})
	err = repo.DeleteFood(1)
	assert.NoError(t, err)
}

func TestChangeFoodCondition(t *testing.T) {
	repo := initTest(t)

	var err error
	err = repo.ChangeCondition(1, entity.Cooked)
	assert.Error(t, err)

	repo.db.Create(&entity.Food{ID: 1, Condition: entity.Raw})
	err = repo.ChangeCondition(1, entity.Cooked)
	assert.NoError(t, err)
}

func TestStringToFoodCondition(t *testing.T) {
	t.Parallel()

	var err error
	var cond entity.FoodCondition
	cond, err = StringToFoodCondition("raw")
	assert.NoError(t, err)
	assert.Equal(t, entity.Raw, cond)

	_, err = StringToFoodCondition("bad input")
	assert.Error(t, err)
}
