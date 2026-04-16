package templates

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

type TemplateTestSuite struct {
	suite.Suite
	ctx        context.Context
	db         *gorm.DB
	tx         *gorm.DB
	s          *TemplateStore
	piid       uuid.UUID
	ingIds     [2]uint
	templateID uint
}

func (s *TemplateTestSuite) SetupSuite() {
	s.piid = uuid.FromStringOrNil("cf0d4408-8db5-4572-b5d9-4ed873d1341f")
	s.ctx = context.WithValue(context.Background(), contextKeys.Piid, s.piid)
	sqlDb, err := et.NewTestDatabase(s.ctx, "hista_db")
	s.Require().NoError(err)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}), &gorm.Config{TranslateError: true})
	s.Require().NoError(err)
	s.db = db
}

func (suite *TemplateTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.Require().NoError(suite.tx.Error)
	suite.s = NewTemplateStore(suite.tx)

	suite.ingIds = [2]uint{0, 0}
	suite.templateID = 0
}

func (s *TemplateTestSuite) TearDownTest() {
	err := s.tx.Rollback().Error
	s.Require().NoError(err)
}

func TestTemplateTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateTestSuite))
}

func (s *TemplateTestSuite) createTemplate() {
	s.createIngredients()

	ing1ID := s.ingIds[0]
	ing2ID := s.ingIds[1]

	var items = []entity.TemplateItem{
		{Condition: entity.Cooked, IngredientID: ing1ID},
		{Condition: entity.Raw, IngredientID: ing2ID},
	}
	id, err := s.s.Create(s.ctx, "template", items)
	s.Require().NoError(err)
	s.templateID = id
}

func (s *TemplateTestSuite) createIngredients() {
	var ing1 = entity.Ingredient{Name: "name 1"}
	err := generic_queries.Create(s.ctx, s.tx, &ing1)
	s.Require().NoError(err)
	s.ingIds[0] = ing1.ID

	var ing2 = entity.Ingredient{Name: "name 2"}
	err = generic_queries.Create(s.ctx, s.tx, &ing2)
	s.Require().NoError(err)
	s.ingIds[1] = ing2.ID
}

func (s *TemplateTestSuite) getTemplateById(id uint) entity.Template {
	template, err := gorm.G[entity.Template](s.tx).
		Where("id = ?", id).
		Preload("Items", nil).
		First(s.ctx)
	s.Require().NoError(err)
	return template
}

func count[T generic_queries.Piider](s *TemplateTestSuite) int64 {
	no, err := gorm.G[T](s.tx).Count(s.ctx, "*")
	s.Require().NoError(err)
	return no
}
