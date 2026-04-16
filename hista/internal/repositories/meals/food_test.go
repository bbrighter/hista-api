package meals

import (
	"testing"

	"encore.app/hista/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func (s *MealRepoTestSuite) TestGetFood() {
	foods, err := s.repo.ListFoods(s.ctx, 10)
	s.NoError(err)
	s.Len(foods, 0)

	id := s.createIngredient()
	s.createFood(id)

	foods, err = s.repo.ListFoods(s.ctx, s.mealId)
	s.NoError(err)
	s.Len(foods, 1)
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

func (s *MealRepoTestSuite) TestBatchCreateFood() {
	id1 := s.createIngredientWithProps("name 1", nil)
	id2 := s.createIngredientWithProps("name 2", nil)
	s.createMeal()

	var foods = []entity.Food{
		{IngredientID: id1, MealID: s.mealId, Condition: entity.Cooked},
		{IngredientID: id2, MealID: s.mealId, Condition: entity.Raw},
	}
	results, err := s.repo.BatchCreateFoods(s.ctx, foods)

	s.NoError(err)
	s.Len(results, 2)
	for _, r := range results {
		s.NotEqualValues(r.ID, 0)
		s.Equal(s.piid, r.PIID)
		s.Equal(s.piid, r.IngredientPIID)
		s.NotEqualValues(0, r.IngredientID)
		s.EqualValues(0, r.Ingredient.ID, "does not return ingredients")
		s.Equal("", r.Ingredient.Name, "does not return ingredients")
	}
}
