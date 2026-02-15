package meals

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/contextKeys"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MealRepoTestSuite struct {
	suite.Suite
	ctx  context.Context
	db   *gorm.DB
	repo *MealRepository
}

const GUID_STR = "cf0d4408-8db5-4572-b5d9-4ed873d1341f"

var GUID = uuid.FromStringOrNil(GUID_STR)

func (suite *MealRepoTestSuite) SetupSuite() {
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
	suite.repo = NewMealRepository(db)
	suite.fillWithData()
}

func (s *MealRepoTestSuite) fillWithData() {
	var ing1 = entity.Ingredient{ID: 1, PIID: GUID, Name: "Ing1", IsArchived: false}
	var ing2 = entity.Ingredient{ID: 2, PIID: GUID, Name: "Ing2", IsArchived: false}
	var ingArchived = entity.Ingredient{ID: 3, PIID: GUID, Name: "ArchivedIng", IsArchived: true}
	var meal1 = entity.Meal{ID: 1, PIID: GUID, Date: time.Date(2025, 6, 6, 6, 0, 0, 0, time.UTC), Freshness: entity.Fresh, StressLevel: 0, IsAlone: true,
		Foods: []entity.Food{
			{PIID: GUID, Condition: entity.Cooked, Ingredient: ing1},
			{PIID: GUID, Condition: entity.Raw, Ingredient: ing2},
		}}
	var meal2 = entity.Meal{ID: 2, PIID: GUID, Date: time.Date(2022, 6, 6, 6, 0, 0, 0, time.UTC), Freshness: entity.Fresh, StressLevel: 0, IsAlone: true,
		Foods: []entity.Food{
			{PIID: GUID, Condition: entity.Cooked, Ingredient: ingArchived},
		}}

	err := s.db.CreateInBatches(&entity.Meals{&meal1, &meal2}, 10).Error
	s.Require().NoError(err)

}

func TestMealRepoTestSuite(t *testing.T) {
	suite.Run(t, new(MealRepoTestSuite))
}
