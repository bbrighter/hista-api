package statistics

import (
	"testing"
	"time"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initTest(t *testing.T) *StatisticsRepo {
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
	return &StatisticsRepo{db: db}
}

func TestFindFoodForSymptoms(t *testing.T) {
	t.Skip() // Doesn't work for SQLite, has to be tested in integration tests
	repo := initTest(t)

	resp, err := repo.FindFoodForSymptoms(time.Now(), time.Now().Add(time.Minute), []uint{1})
	assert.NoError(t, err)
	assert.Len(t, resp, 0)
}

func TestCoundFoods(t *testing.T) {
	repo := initTest(t)

	results := repo.CountFoods([]uint{1})
	assert.Len(t, results, 0)

	// TODO: test with data filled
}
