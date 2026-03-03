package meals

import (
	"testing"

	"encore.app/hista/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGetFood(t *testing.T) {
	repo, ctx := initTest(t)

	foods, err := repo.ListFoods(ctx, 10)
	assert.NoError(t, err)
	assert.Len(t, foods, 0)

	repo.db.Create(&entity.Food{ID: 1, MealID: 10, PIID: GUID})

	foods, err = repo.ListFoods(ctx, 10)
	assert.NoError(t, err)
	assert.Len(t, foods, 1)
}

func TestCreateFoodByName(t *testing.T) {
	repo, ctx := initTest(t)

	var err error
	var meal = entity.Meal{ID: 1, PIID: GUID}
	var food = &entity.Food{MealID: 1, PIID: GUID}
	repo.db.Create(&meal)
	err = repo.CreateFoodByName(ctx, food, "New name")
	assert.NoError(t, err)
	assert.EqualValues(t, 1, food.ID)
	assert.EqualValues(t, 1, food.Ingredient.ID)
}

func TestCreateFoodById(t *testing.T) {
	repo, ctx := initTest(t)

	var err error
	var food = &entity.Food{
		MealID:       1,
		IngredientID: 10,
		Condition:    entity.Cooked,
		PIID:         GUID,
	}
	err = repo.CreateFoodByID(ctx, food)
	assert.Error(t, err)

	var meal = entity.Meal{ID: 1, PIID: GUID}
	repo.db.Create(&meal)
	repo.db.Create(&entity.Ingredient{ID: 2})
	food.IngredientID = 2
	err = repo.CreateFoodByID(ctx, food)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, food.ID)
}

func TestDeleteFood(t *testing.T) {
	repo, ctx := initTest(t)

	var err error
	err = repo.DeleteFood(ctx, 1)
	assert.Error(t, err)

	var food = entity.Food{ID: 1, Ingredient: entity.Ingredient{ID: 10, PIID: GUID}, PIID: GUID}
	repo.db.Create(&food)
	err = repo.DeleteFood(ctx, 1)
	assert.NoError(t, err)

	count, _ := gorm.G[entity.Ingredient](repo.db).Count(ctx, "*")
	assert.EqualValues(t, 0, count)
}

func TestDeleteFoodKeepsUsedIngredients(t *testing.T) {
	repo, ctx := initTest(t)

	var food1 = entity.Food{ID: 1, Ingredient: entity.Ingredient{ID: 10, PIID: GUID}, PIID: GUID}
	var food2 = entity.Food{ID: 2, Ingredient: entity.Ingredient{ID: 10, PIID: GUID}, PIID: GUID}
	err := gorm.G[[]entity.Food](repo.db).Create(ctx, &[]entity.Food{food1, food2})
	require.NoError(t, err)

	err = repo.DeleteFood(ctx, 1)

	assert.NoError(t, err)
	count, _ := gorm.G[entity.Ingredient](repo.db).Where("id = 10").Count(ctx, "*")
	assert.EqualValues(t, 1, count)
}
