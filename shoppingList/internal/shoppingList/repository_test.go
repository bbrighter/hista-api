package shoppinglist

import (
	"context"
	"testing"

	"encore.app/shared/contextKeys"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repoTestSuite struct {
	suite.Suite
	repo *ShoppingListRepo
	ctx  context.Context
	piid uuid.UUID
	db   *gorm.DB
	tx   *gorm.DB
}

func (s *repoTestSuite) SetupSuite() {
	const GUID_STR = "cf0d4408-8db5-4572-b5d9-4ed873d1341f"
	s.piid = uuid.FromStringOrNil(GUID_STR)
	s.ctx = context.WithValue(context.Background(), contextKeys.Piid, s.piid)
	sqlDb, err := et.NewTestDatabase(s.ctx, "shopping_list")
	s.Require().NoError(err)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}), &gorm.Config{TranslateError: true})
	s.db = db
	s.Require().NoError(err)
}
func (s *repoTestSuite) Debug() {
	s.tx = s.tx.Debug()
}

func (s *repoTestSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.Require().NoError(s.tx.Error)
	s.repo = NewShoppingListRepo(s.tx)
}

func (s *repoTestSuite) TearDownTest() {
	err := s.tx.Rollback().Error
	s.Require().NoError(err)
}

func TestShoppingListRepo(t *testing.T) {
	suite.Run(t, new(repoTestSuite))
}

func (s *repoTestSuite) TestCreateListEnforcesUniqueness() {
	var list = &List{}

	err := s.repo.CreateShoppingList(s.ctx, list)
	s.NoError(err)

	err = s.repo.CreateShoppingList(s.ctx, list)
	s.Error(err)
	s.ErrorIs(err, gorm.ErrDuplicatedKey)
}

func (s *repoTestSuite) TestCreateListEnforcesUniquenessAfterDeletionOk() {
	var list = &List{}

	err := s.repo.CreateShoppingList(s.ctx, list)
	s.NoError(err)

	err = s.repo.DeleteShoppingList(s.ctx, list.ID)
	s.NoError(err)

	var newList = &List{}
	err = s.repo.CreateShoppingList(s.ctx, newList)
	s.NoError(err)
}
