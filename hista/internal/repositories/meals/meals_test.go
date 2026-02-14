package meals

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/contextKeys"
	"encore.app/shared/generic_queries"
	"encore.dev/types/uuid"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func initTest(t *testing.T) (*MealRepository, context.Context) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(
		&entity.Meal{},
		&entity.Ingredient{},
		&entity.Food{},
	)
	require.NoError(t, err)

	ctx := context.WithValue(t.Context(), contextKeys.Piid, GUID)
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
	tests := map[string]struct {
		expectedLen  int
		createMeal   bool
		useWrongPiid bool
	}{
		"zero meals": {expectedLen: 0},
		"one meal":   {createMeal: true, expectedLen: 1},
		"wrong piid": {createMeal: true, useWrongPiid: true, expectedLen: 0},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, ctx := initTest(t)
			createContext := ctx
			if test.useWrongPiid {
				newGuid, _ := uuid.NewV4()
				createContext = context.WithValue(createContext, contextKeys.Piid, newGuid)
			}
			if test.createMeal {
				err := generic_queries.Create(createContext, repo.db, &entity.Meal{})
				require.NoError(t, err)
			}
			meals, err := repo.ListMeals(ctx)
			assert.NoError(t, err)
			assert.Len(t, meals, test.expectedLen)
		})
	}
}

func TestGetMeal(t *testing.T) {
	tests := map[string]struct {
		useNonExistingId bool
		useWrongPiid     bool
		expectError      error
	}{
		"ok":         {},
		"not found":  {useNonExistingId: true, expectError: gorm.ErrRecordNotFound},
		"wrong piid": {useWrongPiid: true, expectError: gorm.ErrRecordNotFound},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, ctx := initTest(t)
			createContext := ctx
			if test.useWrongPiid {
				newGuid, _ := uuid.NewV4()
				createContext = context.WithValue(createContext, contextKeys.Piid, newGuid)
			}
			meal1 := entity.Meal{Freshness: entity.Fresh}
			err := generic_queries.Create(createContext, repo.db, &meal1)
			require.NoError(t, err)
			var ing = entity.Ingredient{Name: "name"}
			err = generic_queries.Create(createContext, repo.db, &ing)
			require.NoError(t, err)
			var food = entity.Food{MealID: meal1.ID, IngredientID: ing.ID}
			err = generic_queries.Create(createContext, repo.db, &food)
			require.NoError(t, err)
			meal2 := entity.Meal{Freshness: entity.Older}
			err = generic_queries.Create(createContext, repo.db, &meal2)
			require.NoError(t, err)

			var id1 uint = 1000
			var id2 uint = 1000
			if !test.useNonExistingId {
				id1 = meal1.ID
				id2 = meal2.ID
			}

			meal, err := repo.GetMeal(ctx, id1)
			if test.expectError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, entity.Fresh, meal.Freshness)
				assert.Len(t, meal.Foods, 1, "Foods are returned via preload")
				otherMeal, _ := repo.GetMeal(ctx, id2)
				assert.Equal(t, entity.Older, otherMeal.Freshness, "The correct meal is returned, not just any")
			}

		})
	}
}

func TestDeleteMeal(t *testing.T) {
	tests := map[string]struct {
		useExistingId bool
		expectedError error
	}{
		"ok":        {useExistingId: true},
		"not found": {useExistingId: false, expectedError: gorm.ErrRecordNotFound},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, ctx := initTest(t)
			meal, ing, food := createTestMeal(t, ctx, repo.db)
			var id uint = 1000
			if test.useExistingId {
				id = meal.ID
			}
			err := repo.DeleteMeal(ctx, id)
			if test.expectedError != nil {
				assert.ErrorIs(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
				foodCount, _ := generic_queries.Count[*entity.Food](ctx, repo.db, food.ID)
				assert.EqualValues(t, 0, foodCount)
				ingredientCount, _ := generic_queries.Count[*entity.Ingredient](ctx, repo.db, ing.ID)
				assert.EqualValues(t, 0, ingredientCount)
			}
		})
	}
}

func TestPatchMeal(t *testing.T) {
	var newDate time.Time = time.Date(2020, 5, 5, 5, 5, 0, 0, time.UTC)
	var newFreshness entity.Freshness = entity.Older // Default = entity.Fresh
	var newStressLevel uint8 = 3                     // Default = 0
	var newIsAlone bool = true                       // Default = false
	tests := map[string]struct {
		useExistingId    bool
		useWrongPiid     bool
		patchDate        *time.Time
		patchFreshness   *entity.Freshness
		patchStressLevel *uint8
		patchIsAlone     *bool
		expectedError    error
	}{
		"ok, patch all":          {useExistingId: true, patchDate: &newDate, patchFreshness: &newFreshness, patchStressLevel: &newStressLevel, patchIsAlone: &newIsAlone},
		"ok, patch date":         {useExistingId: true, patchDate: &newDate},
		"ok, patch freshness":    {useExistingId: true, patchFreshness: &newFreshness},
		"ok, patch stress level": {useExistingId: true, patchStressLevel: &newStressLevel},
		"ok, patch isAlone":      {useExistingId: true, patchIsAlone: &newIsAlone},
		"not found":              {useExistingId: false, patchDate: &newDate, expectedError: gorm.ErrRecordNotFound},
		"no change params":       {useExistingId: true, expectedError: ErrNoParameters},
		"wrong piid":             {useExistingId: true, patchDate: &newDate, useWrongPiid: true, expectedError: gorm.ErrRecordNotFound},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, ctx := initTest(t)
			meal, _, _ := createTestMeal(t, ctx, repo.db)

			var id uint = 1000
			if test.useExistingId {
				id = meal.ID
			}
			patchCtx := ctx
			if test.useWrongPiid {
				newGuid, _ := uuid.NewV4()
				patchCtx = context.WithValue(patchCtx, contextKeys.Piid, newGuid)
			}
			err := repo.PatchMeal(patchCtx, id, test.patchDate, test.patchFreshness, test.patchStressLevel, test.patchIsAlone)
			if test.expectedError != nil {
				assert.ErrorIs(t, err, test.expectedError)
				return
			} else {
				assert.NoError(t, err)
			}

			dbMeal, err := generic_queries.First[*entity.Meal](ctx, repo.db, meal.ID)
			require.NoError(t, err)
			assert.Equal(t, test.patchDate != nil, dbMeal.Date.Equal(newDate))
			assert.Equal(t, test.patchFreshness != nil, dbMeal.Freshness == newFreshness)
			assert.Equal(t, test.patchIsAlone != nil, dbMeal.IsAlone == newIsAlone)
			assert.Equal(t, test.patchStressLevel != nil, dbMeal.StressLevel == newStressLevel)
		})
	}
}

func createTestMeal(t *testing.T, ctx context.Context, db *gorm.DB) (entity.Meal, entity.Ingredient, entity.Food) {
	var meal = entity.Meal{Freshness: entity.Fresh}
	err := generic_queries.Create(ctx, db, &meal)
	require.NoError(t, err)
	var ing = entity.Ingredient{Name: "name"}
	err = generic_queries.Create(ctx, db, &ing)
	require.NoError(t, err)
	var food = entity.Food{MealID: meal.ID, IngredientID: ing.ID}
	err = generic_queries.Create(ctx, db, &food)
	require.NoError(t, err)
	return meal, ing, food
}
