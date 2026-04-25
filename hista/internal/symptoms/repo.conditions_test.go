package symptoms

import (
	"context"
	"testing"
	"time"

	"encore.app/shared/contextKeys"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type symSuite struct {
	suite.Suite
	db  *gorm.DB
	tx  *gorm.DB
	ctx context.Context
	con *ConditionRepo
	sym *SymptomRepo
}

const GUID_STR = "cf0d4408-8db5-4572-b5d9-4ed873d1341f"

func (s *symSuite) SetupSuite() {
	piid := uuid.FromStringOrNil(GUID_STR)
	s.ctx = context.WithValue(context.Background(), contextKeys.Piid, piid)
	sqlDb, err := et.NewTestDatabase(s.ctx, "hista_db")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}), &gorm.Config{TranslateError: true})
	s.db = db
	if err != nil {
		panic(err)
	}
}

func (s *symSuite) Debug() {
	s.tx = s.tx.Debug()
}

func (s *symSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.Require().NoError(s.tx.Error)
	s.con = NewConditionRepo(s.tx)
	s.sym = NewSymptomRepo(s.tx)
}

func (s *symSuite) TearDownTest() {
	err := s.tx.Rollback().Error
	s.Require().NoError(err)
}

func TestSymptoms(t *testing.T) {
	suite.Run(t, new(symSuite))
}

func (s *symSuite) createTestEvent() uint {
	var event = &ConditionEvent{Date: time.Now()}
	err := s.con.CreateConditionEvent(s.ctx, event)
	s.Require().NoError(err)
	return event.ID
}

func (s *symSuite) createTestCategory() uint {
	var cat = &SymptomCategory{Name: "Category"}
	err := s.sym.CreateSymptomCategory(s.ctx, cat)
	s.Require().NoError(err)
	return cat.ID
}

func (s *symSuite) addConditionToEvent(eventId uint) uint {
	catId := s.createTestCategory()
	var symptom = Symptom{Name: "Symptom", SymptomCategoryID: catId}
	var condition = &Condition{Severity: 1, Symptom: symptom, ConditionEventID: eventId}
	err := s.con.CreateCondition(s.ctx, condition)
	s.Require().NoError(err)
	return condition.ID
}

func (s *symSuite) TestListConditionEventsAndDependencies() {
	events, err := s.con.ListConditionEventsAndDependencies(s.ctx)
	s.NoError(err)
	s.Len(events, 0)

	eventId := s.createTestEvent()
	s.addConditionToEvent(eventId)

	events, err = s.con.ListConditionEventsAndDependencies(s.ctx)
	s.NoError(err)
	s.Len(events, 1)
	event := events[0]
	s.Len(event.Conditions, 1)
	cond := event.Conditions[0]
	s.Equal("Symptom", cond.Symptom.Name)
}

func (s *symSuite) TestListConditions() {
	eventId := s.createTestEvent()
	s.addConditionToEvent(eventId)

	conditions, err := s.con.ListConditions(s.ctx, eventId)
	s.NoError(err)
	s.Len(conditions, 1)
}

func (s *symSuite) TestListConditions_notFound() {
	conditions, err := s.con.ListConditions(s.ctx, 100)
	s.NoError(err)
	s.Len(conditions, 0)
}

func (s *symSuite) TestListSymptomCategoriesAndSymptoms() {
	eventId := s.createTestEvent()
	s.addConditionToEvent(eventId)

	cats, err := s.sym.ListSymptomCategoriesAndSymptoms(s.ctx)
	s.NoError(err)
	s.Len(cats, 1)
	cat := cats[0]
	s.Equal("Category", cat.Name)
	s.Len(cat.Symptoms, 1)
	sym := cat.Symptoms[0]
	s.Equal("Symptom", sym.Name)
}
