package meals

import (
	"context"
	"testing"

	"encore.app/hista/entity"
	"encore.app/shared/contextKeys"
	"encore.app/shared/generic_queries"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MealRepoTestSuite struct {
	suite.Suite
	ctx          context.Context
	db           *gorm.DB // DB connection
	tx           *gorm.DB // Transaction
	repo         *MealRepository
	piid         uuid.UUID
	mealId       uint
	ingredientId uint
}

const GUID_STR = "cf0d4408-8db5-4572-b5d9-4ed873d1341f"

var GUID = uuid.FromStringOrNil(GUID_STR)

func (suite *MealRepoTestSuite) SetupSuite() {
	suite.piid = uuid.FromStringOrNil(GUID_STR)
	suite.ctx = context.WithValue(context.Background(), contextKeys.Piid, suite.piid)
	sqlDb, err := et.NewTestDatabase(suite.ctx, "hista_db")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}), &gorm.Config{TranslateError: true})
	suite.db = db
	if err != nil {
		panic(err)
	}
}

func (suite *MealRepoTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.Require().NoError(suite.tx.Error)
	suite.repo = NewMealRepository(suite.tx)
}

func (s *MealRepoTestSuite) TearDownTest() {
	err := s.tx.Rollback().Error
	s.Require().NoError(err)
}

func (s *MealRepoTestSuite) createIngredient() uint {
	return s.createIngredientWithProps("ingredient", nil)
}

func (s *MealRepoTestSuite) createIngredientWithProps(name string, protein *float32) uint {
	ing := entity.Ingredient{Name: name, Nutrition: entity.Nutrition{Protein: protein}}
	err := generic_queries.Create(s.ctx, s.tx, &ing)
	s.Require().NoError(err)
	s.ingredientId = ing.ID
	return ing.ID
}

func (s *MealRepoTestSuite) createFood(ingID uint) uint {
	s.createMeal()

	amount := 20

	food := entity.Food{
		Condition:  entity.Cooked,
		MealID:     s.mealId,
		Amount:     &amount,
		Ingredient: entity.Ingredient{ID: ingID, PIID: s.piid},
	}
	err := generic_queries.Create(s.ctx, s.tx, &food)
	s.Require().NoError(err)
	return food.ID
}

func (s *MealRepoTestSuite) createMeal() {
	meal := entity.Meal{}
	err := generic_queries.Create(s.ctx, s.tx, &meal)
	s.Require().NoError(err)
	s.mealId = meal.ID
}

func TestMealRepoTestSuite(t *testing.T) {
	suite.Run(t, new(MealRepoTestSuite))
}
