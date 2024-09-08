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
	var meal = entity.Meal{ID: 1}
	var food = &entity.Food{MealID: 1}
	repo.db.Create(&meal)
	err = repo.CreateFoodByName(food, "New name")
	assert.NoError(t, err)
	assert.EqualValues(t, 1, food.ID)
	assert.EqualValues(t, 1, food.Ingredient.ID)
}

func TestCreateFoodById(t *testing.T) {
	repo := initTest(t)

	var err error
	var food = &entity.Food{
		MealID:       1,
		IngredientID: 10,
		Condition:    entity.Cooked,
	}
	err = repo.CreateFoodByID(food)
	assert.Error(t, err)

	var meal = entity.Meal{ID: 1}
	repo.db.Create(&meal)
	repo.db.Create(&entity.Ingredient{ID: 2})
	food.IngredientID = 2
	err = repo.CreateFoodByID(food)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, food.ID)
}

func TestDeleteFood(t *testing.T) {
	repo := initTest(t)

	var err error
	err = repo.DeleteFood(1)
	assert.Error(t, err)

	var food = entity.Food{ID: 1, Ingredient: entity.Ingredient{ID: 10}}
	repo.db.Create(&food)
	err = repo.DeleteFood(1)
	assert.NoError(t, err)
}

func TestChangeFoodCondition(t *testing.T) {
	repo := initTest(t)

	var err error
	var condition = entity.Raw
	err = repo.ChangeCondition(1, condition)
	assert.Error(t, err)

	var food = entity.Food{ID: 1, Condition: entity.Cooked}
	repo.db.Create(&food)
	err = repo.ChangeCondition(1, condition)
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
