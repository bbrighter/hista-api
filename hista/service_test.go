package hista

import (
	"context"
	"fmt"
	"testing"
	"time"

	"encore.app/hista/internal/dwd"
	"encore.app/shared/contextKeys"
	"encore.dev/beta/errs"
	"encore.dev/et"
	"encore.dev/types/option"
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

type dwdClientMock struct{}

func (d dwdClientMock) GetKarlsruheData() (dwd.DWDPollen, time.Time, error) {
	dwd := dwd.DWDPollen{
		Hasel:    dwd.DWDPollenIntensity{Today: "0"},
		Esche:    dwd.DWDPollenIntensity{Today: "0"},
		Graeser:  dwd.DWDPollenIntensity{Today: "1"},
		Ambrosia: dwd.DWDPollenIntensity{Today: "2"},
		Erle:     dwd.DWDPollenIntensity{Today: "1-2"},
		Roggen:   dwd.DWDPollenIntensity{Today: "0-1"},
		Birke:    dwd.DWDPollenIntensity{Today: "3"},
		Beifuss:  dwd.DWDPollenIntensity{Today: "2-3"},
	}
	updatedAt := time.Date(2024, 5, 31, 11, 0, 0, 0, time.UTC)
	return dwd, updatedAt, nil
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
	}), &gorm.Config{TranslateError: true})
	s.Require().NoError(err)
	s.db = db
	s.service = *initServiceWithDb(s.db, new(dwdClientMock))
}

func (s *ApiTestSuite) Debug() {
	s.service = *initServiceWithDb(s.db.Debug(), new(dwdClientMock))
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
		"templates", "template_items",
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
	meal, err := s.service.PostMeal(s.ctx, s.piid, PostMealParams{Date: time.Now()})
	s.Require().NoError(err)
	return meal.ID
}

func (s *ApiTestSuite) createTestFood() (foodId uint, ingredientId uint) {
	mealId := s.createTestMeal()
	foodResp, err := s.service.PostFood(s.ctx, s.piid, mealId, FoodParams{IngredientName: "ingredient", IngredientID: 0})
	s.Require().NoError(err)
	return foodResp.Food.ID, foodResp.Food.IngredientId
}

func (s *ApiTestSuite) createTestCondition() (condId uint, symptomId uint, catId uint) {
	id := s.createTestEvent()
	catResp, err := s.service.PostSymptomCategory(s.ctx, s.piid, PostSymptomCategoryRequest{Name: "cat"})
	s.Require().NoError(err)
	var name string = "name"
	resp, err := s.service.PostCondition(s.ctx, s.piid, id, ConditionRequestParams{SymptomName: option.Some(name), CategoryID: option.Some(catResp.ID)})
	s.Require().NoError(err)

	if resp.Symptoms.IsSome() {
		symptoms := resp.Symptoms.MustGet()
		symptomId = symptoms.Categories[0].Symptoms[0].ID
		catId = symptoms.Categories[0].ID
	}

	return resp.Condition.ID, symptomId, catId
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
	err := s.service.UpdatePollen(s.ctx)
	s.Require().NoError(err)
}

func (s *ApiTestSuite) createTestIntake() {
	id, err := s.service.meds.CreateMedicine(s.ctx, "Medicine No. 1")
	s.Require().NoError(err)
	err = s.service.meds.IncrementIntake(s.ctx, id)
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

func (s *ApiTestSuite) createTestTemplate() uint {
	_, ingId := s.createTestFood()
	template, err := s.service.PostTemplate(s.ctx, s.piid, TemplateParams{Name: "Template", Items: []TemplateItemParams{
		{IngredientId: ingId, Condition: "cooked"},
	}})
	s.Require().NoError(err)
	return template.ID
}
