package hista

import (
	"context"
	"fmt"
	"testing"
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/contextKeys"
	"encore.dev/beta/errs"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ApiTestSuite struct {
	suite.Suite
	piid    uuid.UUID
	service Service
	ctx     context.Context
	db      *gorm.DB
}

func (s *ApiTestSuite) SetupSuite() {
	guid, err := uuid.FromString("0c5e945e-ef6c-4934-91ff-702d94e2e7a8")
	s.Require().NoError(err)
	s.piid = guid
	s.ctx = context.WithValue(context.Background(), contextKeys.Piid, guid)
	sqlDb, err := et.NewTestDatabase(s.ctx, "hista_db")
	s.Require().NoError(err)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}))
	s.Require().NoError(err)
	s.db = db
	s.service = *initServiceWithDb(s.db)
}

func (s *ApiTestSuite) Debug() {
	s.service = *initServiceWithDb(s.db.Debug())
}

func (s *ApiTestSuite) cleanTables() {
	tables := []string{
		"foods", "ingredients", "meals",
		"conditions", "symptoms", "symptom_categories", "condition_events",
		"notes",
		"pollens", "pollen_events",
		"statuses",
		"headaches",
		"intakes", "medicines",
	}
	for _, table := range tables {
		err := s.db.Exec(fmt.Sprintf(`DELETE FROM "%s"`, table)).Error
		s.Require().NoError(err)
	}
}

func (s *ApiTestSuite) TearDownSubTest() {
	s.cleanTables()
}

func (s *ApiTestSuite) TearDownTest() {
	s.cleanTables()
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func (s *ApiTestSuite) createTestNote() uint {
	resp, err := s.service.PostNote(s.ctx, s.piid)
	s.Require().NoError(err)
	return resp.ID
}

func (s *ApiTestSuite) createTestMeal() uint {
	meal, err := s.service.PostMeal(s.ctx, s.piid, entity.PostMealParams{Date: time.Now()})
	s.Require().NoError(err)
	return meal.ID
}

func (s *ApiTestSuite) createTestFood() (foodId uint, ingredientId uint) {
	meal, err := s.service.PostMeal(s.ctx, s.piid, entity.PostMealParams{Date: time.Now()})
	s.Require().NoError(err)
	foodResp, err := s.service.PostFood(s.ctx, s.piid, meal.ID, FoodParams{IngredientName: "ingredient", IngredientID: 0})
	s.Require().NoError(err)
	return foodResp.Food.ID, foodResp.Food.Ingredient.ID
}

func (s *ApiTestSuite) createTestCondition() (condId uint, symptomId uint, catId uint) {
	id := s.createTestEvent()
	catResp, err := s.service.PostSymptomCategory(s.ctx, s.piid, PostSymptomCategoryRequest{Name: "cat"})
	s.Require().NoError(err)
	var name string = "name"
	resp, err := s.service.PostCondition(s.ctx, s.piid, id, ConditionRequestParams{SymptomName: &name, CategoryID: &catResp.ID})
	s.Require().NoError(err)

	return resp.Condition.ID, resp.Condition.Symptom.ID, resp.Condition.Symptom.CategoryID
}

func (s *ApiTestSuite) createTestEvent() uint {
	resp, err := s.service.CreateConditionEvent(s.ctx, s.piid)
	s.Require().NoError(err)
	return resp.ID
}

func (s *ApiTestSuite) createTestHeadache() uint {
	resp, err := s.service.PostHeadache(s.ctx, s.piid, PostHeadacheParams{Date: time.Now(), Severity: 3})
	s.Require().NoError(err)
	return resp.ID
}

func (s *ApiTestSuite) createTestStatus() uint {
	resp, err := s.service.PostStatus(s.ctx, s.piid, DateParam{Date: time.Now()})
	s.Require().NoError(err)
	return resp.ID
}

func (s *ApiTestSuite) createTestPollen() {
	s.service.pollens.UseTestQuery(s.T())
	err := s.service.UpdatePollen(s.ctx)
	s.Require().NoError(err)
}

func (suite *ApiTestSuite) assertErrCode(err error, expectedCode errs.ErrCode) bool {
	if expectedCode == 0 {
		suite.NoError(err)
		return false
	}
	suite.Require().NotNil(err)
	if err != nil {
		print(err.Error())
	}
	encoreErr, ok := err.(*errs.Error)
	suite.Require().True(ok)
	suite.Equal(encoreErr.Code, expectedCode, encoreErr)
	return true
}
