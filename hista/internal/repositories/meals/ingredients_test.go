package meals

import (
	"testing"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// func TestCreateOrReplaceIngredient(t *testing.T) {
// 	t.Skip()
// 	service := initTest(t)
// 	var ing entity.Ingredient
// 	var err error

// 	ing, err = service.createOrReplaceIngredient("Name")
// 	assert.NoError(t, err)
// 	assert.EqualValues(t, "Name", ing.Name)

// 	// Verify idempotency
// 	var firstId = ing.ID
// 	ing, err = service.createOrReplaceIngredient("Name")
// 	assert.NoError(t, err)
// 	assert.EqualValues(t, firstId, ing.ID)

// 	// Verify idempotency even when trimming
// 	ing, err = service.createOrReplaceIngredient(" Name ")
// 	assert.NoError(t, err)
// 	assert.EqualValues(t, firstId, ing.ID)

// 	// Verify that new items get new names
// 	ing, err = service.createOrReplaceIngredient("Name2")
// 	assert.NoError(t, err)
// 	assert.NotEqualValues(t, firstId, ing.ID)

// }

func TestGetIngredients(t *testing.T) {
	repo, ctx := initTest(t)

	ingredients, err := repo.ListIngredients(ctx)
	assert.NoError(t, err)
	assert.Len(t, ingredients, 0)

	repo.db.Create(&entity.Ingredient{Name: "ingredient", PIID: uuid.FromStringOrNil(GUID_STR)})
	ingredients, err = repo.ListIngredients(ctx)
	assert.NoError(t, err)
	assert.Len(t, ingredients, 1)
}

// func TestDeleteIngredientIfUnused(t *testing.T) {
// 	service := initTest(t)

// 	service.db.Create(&entity.Food{ID: 1, Ingredient: entity.Ingredient{ID: 10}})

// 	var err error
// 	err = deleteIngredientIfUnused(service.db, 10)
// 	assert.NoError(t, err)
// 	// Ingredient is used and should not be removed
// 	var rows int64
// 	var ingredient entity.Ingredient
// 	rows = service.db.Find(&ingredient).RowsAffected
// 	assert.EqualValues(t, 1, rows)

// 	// Ingredient is not used and should be removed
// 	service.db.Create(&entity.Ingredient{ID: 1, Name: "Unused"})
// 	err = deleteIngredientIfUnused(service.db, 1)
// 	assert.NoError(t, err)

// 	rows = service.db.First(entity.Ingredient{ID: 1}).RowsAffected
// 	assert.EqualValues(t, 0, rows)
// }

func (s *MealRepoTestSuite) TestGetIngredients() {
	ing, err := s.repo.ListIngredients(s.ctx)
	s.NoError(err)
	s.Len(ing, 3)
}

func (s *MealRepoTestSuite) TestChangeName() {
	var err error

	err = s.repo.ChangeIngredientName(s.ctx, 1, "new name")
	s.NoError(err)
	ing, err := gorm.G[entity.Ingredient](s.db).Where("id = ?", 1).First(s.ctx)
	s.Require().NoError(err)
	s.Equal("new name", ing.Name)

	err = s.repo.ChangeIngredientName(s.ctx, 100, "new name")
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *MealRepoTestSuite) TestToggleArchived() {
	var err error
	var ing entity.Ingredient

	err = s.repo.ToggleArchived(s.ctx, 1)
	s.NoError(err)
	ing, err = gorm.G[entity.Ingredient](s.db).Where("id = ?", 1).First(s.ctx)
	s.Require().NoError(err)
	s.Equal(true, ing.IsArchived)
	err = s.repo.ToggleArchived(s.ctx, 1)
	ing, err = gorm.G[entity.Ingredient](s.db).Where("id = ?", 1).First(s.ctx)
	s.Require().NoError(err)
	s.Equal(false, ing.IsArchived)

	err = s.repo.ToggleArchived(s.ctx, 100)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *MealRepoTestSuite) TestDeleteIngredient() {
	var err error

	err = s.repo.DeleteIngredient(s.ctx, 1)
	s.Error(err)
	s.ErrorContains(err, "23503")

	ing := entity.Ingredient{ID: 95, Name: "no relations"}
	generic_queries.Create(s.ctx, s.db, &ing)
	err = s.repo.DeleteIngredient(s.ctx, ing.ID)
	s.NoError(err)
}
