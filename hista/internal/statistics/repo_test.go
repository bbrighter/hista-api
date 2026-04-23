package statistics

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/internal/meals"
	"encore.app/hista/internal/symptoms"
	"encore.app/shared/contextKeys"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type StatsRepoTestSuite struct {
	suite.Suite
	ctx   context.Context
	db    *gorm.DB
	stats *statisticsRepo
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
	suite.stats = newStatisticsRepo(db)
	suite.fillWithData()
}

func TestStatisticsRepoTestSuite(t *testing.T) {
	suite.Run(t, new(StatsRepoTestSuite))
}

func (s *StatsRepoTestSuite) fillWithData() {
	mealWithSymptomsDate := time.Date(2022, 1, 9, 12, 0, 0, 0, time.UTC)
	conditionWithin24HDate := mealWithSymptomsDate.Add(time.Hour * 12)
	conditionWihtin72HDate := mealWithSymptomsDate.Add(time.Hour * 52)
	mealWithNoSymptomsDate := mealWithSymptomsDate.Add(time.Hour * 24 * 31)
	cat := symptoms.SymptomCategory{Name: "cat"}
	con1 := symptoms.Condition{
		Severity: symptoms.LowSeverity,
		Symptom:  symptoms.Symptom{ID: 1, Name: "Sym", SymptomCategory: cat},
	}
	con2 := symptoms.Condition{
		Severity: symptoms.LowSeverity,
		Symptom:  symptoms.Symptom{ID: 2, Name: "Sym 2", SymptomCategory: cat},
	}
	eventWithin24H := symptoms.ConditionEvent{
		PIID:       GUID,
		Date:       conditionWithin24HDate,
		Conditions: []symptoms.Condition{con1},
	}
	eventWithin72H := symptoms.ConditionEvent{
		PIID:       GUID,
		Date:       conditionWihtin72HDate,
		Conditions: []symptoms.Condition{con1, con2},
	}
	ing := meals.Ingredient{ID: 1}
	meal1 := meals.Meal{
		ID:    1,
		PIID:  GUID,
		Date:  mealWithSymptomsDate,
		Foods: []meals.Food{{Ingredient: ing}},
	}
	meal2 := meals.Meal{
		ID:    2,
		PIID:  GUID,
		Date:  mealWithNoSymptomsDate,
		Foods: []meals.Food{{Ingredient: ing}},
	}
	err := gorm.G[symptoms.ConditionEvent](s.db).CreateInBatches(s.ctx, &[]symptoms.ConditionEvent{eventWithin24H, eventWithin72H}, 100)
	s.Require().NoError(err)
	err = gorm.G[meals.Meal](s.db).CreateInBatches(s.ctx, &[]meals.Meal{meal1, meal2}, 100)
	s.Require().NoError(err)
}

func (s *StatsRepoTestSuite) TestSymptomsAfterIngredients() {
	validFromDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := map[string]struct {
		fromDate      time.Time
		toDate        time.Time
		ingredientIds uint
		want          []FoodResult
	}{
		"ok": {
			fromDate:      validFromDate,
			toDate:        validFromDate.Add(time.Hour * 24 * 365),
			ingredientIds: 1,
			want: []FoodResult{
				{SymptomID: 1, Hours1: 0, Hours24: 1, Hours72: 1, Severity: 2},
				{SymptomID: 2, Hours1: 0, Hours24: 0, Hours72: 1, Severity: 2},
			},
		},
		"outside date range": {
			fromDate:      validFromDate.Add(time.Hour * 24 * 365),
			toDate:        validFromDate.Add(time.Hour * 24 * 366),
			ingredientIds: 1,
			want:          []FoodResult{},
		},
		"no ingredient": {
			fromDate:      validFromDate,
			toDate:        validFromDate.Add(time.Hour * 24 * 365),
			ingredientIds: 100,
			want:          []FoodResult{},
		},
	}
	for name, tt := range tests {
		s.Run(name, func() {
			got, gotErr := s.stats.SymptomsAfterIngredients(s.ctx, tt.fromDate, tt.toDate, tt.ingredientIds)
			s.NoError(gotErr)
			s.ElementsMatch(got, tt.want)

		})
	}
}

func (s *StatsRepoTestSuite) TestCountMealsWithIngredients() {
	validFromDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := map[string]struct {
		fromDate      time.Time
		toDate        time.Time
		ingredientIds uint
		want          int64
	}{
		"ok": {
			fromDate:      validFromDate,
			toDate:        validFromDate.Add(time.Hour * 24 * 365),
			ingredientIds: 1,
			want:          2,
		},
		"outside date range": {
			fromDate:      validFromDate.Add(time.Hour * 24 * 365),
			toDate:        validFromDate.Add(time.Hour * 24 * 366),
			ingredientIds: 1,
			want:          0,
		},
		"no ingredient": {
			fromDate:      validFromDate,
			toDate:        validFromDate.Add(time.Hour * 24 * 365),
			ingredientIds: 100,
			want:          0,
		},
	}
	for name, tt := range tests {
		s.Run(name, func() {
			got, gotErr := s.stats.CountMealsWithIngredients(s.ctx, tt.fromDate, tt.toDate, tt.ingredientIds)
			s.NoError(gotErr)
			s.EqualValues(got, tt.want)

		})
	}
}
