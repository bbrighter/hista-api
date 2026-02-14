package statistics

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/contextKeys"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type StatsRepoTestSuite struct {
	suite.Suite
	ctx   context.Context
	db    *gorm.DB
	stats *StatisticsRepo
}

const GUID_STR = "cf0d4408-8db5-4572-b5d9-4ed873d1341f"

var GUID = uuid.FromStringOrNil(GUID_STR)

func (suite *StatsRepoTestSuite) SetupSuite() {
	suite.ctx = context.WithValue(context.Background(), contextKeys.Piid, GUID)
	sqlDb, err := et.NewTestDatabase(suite.ctx, "hista_db")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}))
	suite.db = db
	if err != nil {
		panic(err)
	}
	suite.stats = NewStatisticsRepo(db)
	suite.fillWithData()
}

func (s *StatsRepoTestSuite) fillWithData() {
	mealWithSymptomsDate := time.Date(2022, 1, 9, 12, 0, 0, 0, time.UTC)
	conditionWithin24HDate := mealWithSymptomsDate.Add(time.Hour * 12)
	conditionWihtin72HDate := mealWithSymptomsDate.Add(time.Hour * 52)
	mealWithNoSymptomsDate := mealWithSymptomsDate.Add(time.Hour * 24 * 31)
	cat := entity.SymptomCategory{Name: "cat"}
	con1 := entity.Condition{
		Severity: entity.LowSeverity,
		PIID:     GUID,
		Symptom:  entity.Symptom{ID: 1, Name: "Sym", SymptomCategory: cat},
	}
	con2 := entity.Condition{
		Severity: entity.LowSeverity,
		PIID:     GUID,
		Symptom:  entity.Symptom{ID: 2, Name: "Sym 2", SymptomCategory: cat},
	}
	eventWithin24H := entity.ConditionEvent{
		PIID:       GUID,
		Date:       conditionWithin24HDate,
		Conditions: []entity.Condition{con1},
	}
	eventWithin72H := entity.ConditionEvent{
		PIID:       GUID,
		Date:       conditionWihtin72HDate,
		Conditions: []entity.Condition{con1, con2},
	}
	ing := entity.Ingredient{ID: 1}
	meal1 := entity.Meal{
		ID:    1,
		PIID:  GUID,
		Date:  mealWithSymptomsDate,
		Foods: []entity.Food{{PIID: GUID, Ingredient: ing}},
	}
	meal2 := entity.Meal{
		ID:    2,
		PIID:  GUID,
		Date:  mealWithNoSymptomsDate,
		Foods: []entity.Food{{PIID: GUID, Ingredient: ing}},
	}
	err := gorm.G[entity.ConditionEvent](s.db).CreateInBatches(s.ctx, &[]entity.ConditionEvent{eventWithin24H, eventWithin72H}, 100)
	s.Require().NoError(err)
	err = gorm.G[entity.Meal](s.db).CreateInBatches(s.ctx, &[]entity.Meal{meal1, meal2}, 100)
	s.Require().NoError(err)
}

// func (suite *StatsRepoTestSuite) SetupSubTest() {
// 	suite.stats = NewStatisticsRepo(suite.db)
// }

// func (suite *StatsRepoTestSuite) TearDownSubTest() {
// 	var err error
// 	tx := suite.db
// 	var tables = []string{
// 		"meals",
// 		"foods",
// 		"ingredients",
// 		"condition_events",
// 		"conditions",
// 		"symptoms",
// 	}
// 	for _, t := range tables {
// 		err = tx.Exec(`DELETE FROM ` + t).Error
// 		suite.Require().NoError(err)
// 	}
// }

func (s *StatsRepoTestSuite) AssertPostgresError(err error, code string) {
	pgErr, ok := err.(*pgconn.PgError)
	s.True(ok, "expected Postgres error")
	s.Equal(code, pgErr.Code)
}

func TestStatisticsRepoTestSuite(t *testing.T) {
	suite.Run(t, new(StatsRepoTestSuite))
}
