package statistics

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/contextKeys"
	"encore.dev/types/uuid"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var GUID = uuid.FromStringOrNil("5f0347ae-38e8-49f4-9707-c3f4500e1768")

func initTest(t *testing.T) (*StatisticsRepo, context.Context) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(
		&entity.ConditionEvent{},
		&entity.Condition{},
		&entity.Symptom{},
		&entity.SymptomCategories{},
		&entity.Meal{},
		&entity.Food{},
		&entity.Ingredient{},
	)
	assert.NoError(t, err)

	ctx := context.WithValue(t.Context(), contextKeys.Piid, GUID)
	return &StatisticsRepo{db: db}, ctx
}

func TestFindFoodForSymptoms(t *testing.T) {
	t.Skip() // Doesn't work for SQLite, has to be tested in integration tests
	repo, ctx := initTest(t)

	resp, err := repo.FindFoodForSymptoms(ctx, time.Now(), time.Now().Add(time.Minute), []uint{1})
	assert.NoError(t, err)
	assert.Len(t, resp, 0)
}

func TestCoundFoods(t *testing.T) {
	repo, ctx := initTest(t)

	results, err := repo.CountFoods(ctx, []uint{1})
	assert.NoError(t, err)
	assert.Len(t, results, 0)

	// TODO: test with data filled
}
