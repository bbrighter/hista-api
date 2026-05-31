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

func (s *repoTestSuite) createProduct(name string) uint {
	var product = Product{Name: name}
	err := s.repo.UpsertProduct(s.ctx, &product)
	s.Require().NoError(err)
	return product.ID
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

func (s *repoTestSuite) TestDeleteShoppingListKeepsChildren() {
	var list = &List{}
	err := s.repo.CreateShoppingList(s.ctx, list)
	s.NoError(err)
	var item = &Item{ListId: list.ID, Product: Product{Name: "Name"}}
	err = s.repo.CreateItem(s.ctx, item)
	s.NoError(err)

	err = s.repo.DeleteShoppingList(s.ctx, list.ID)
	s.NoError(err)

	count, err := gorm.G[Item](s.tx).Where("list_id = ?", list.ID).Count(s.ctx, "*")
	s.NoError(err)
	s.EqualValues(1, count)
}

func (s *repoTestSuite) TestUpdateProductOk() {
	id := s.createProduct("Name")

	err := s.repo.UpdateProduct(
		s.ctx,
		id,
		map[string]any{
			"archived": true,
			"name":     "new name",
		})
	s.NoError(err)
}
func (s *repoTestSuite) TestUpsertProductNew() {
	var product = Product{Name: "name"}

	err := s.repo.UpsertProduct(s.ctx, &product)
	s.NoError(err)
	s.False(product.Archived)
	s.Equal("name", product.Name)

	prods, err := s.repo.ListProducts(s.ctx)
	s.Require().NoError(err)
	s.Len(prods, 1)
}

func (s *repoTestSuite) TestUpsertProductWithSameNameUnarchives() {
	id := s.createProduct("name")
	err := s.repo.UpdateProduct(s.ctx, id, map[string]any{"archived": true})
	s.Require().NoError(err)

	var product = Product{Name: "name"}

	err = s.repo.UpsertProduct(s.ctx, &product)
	s.NoError(err)
	s.False(product.Archived)
	s.Equal("name", product.Name)

	prods, err := s.repo.ListProducts(s.ctx)
	s.Require().NoError(err)
	s.Len(prods, 1)
	s.False(prods[0].Archived)
}

func (s *repoTestSuite) TestCreateProductWithOtherName() {
	id := s.createProduct("name")

	var product = Product{Name: "other name"}

	err := s.repo.UpsertProduct(s.ctx, &product)
	s.NoError(err)
	s.False(product.Archived)
	s.Equal("other name", product.Name)
	s.NotEqual(id, product.ID)

	prods, err := s.repo.ListProducts(s.ctx)
	s.Require().NoError(err)
	s.Len(prods, 2)
}

func (s *repoTestSuite) TestCreateProductWithSameNameAndOtherPiid() {
	id := s.createProduct("name")

	ctx := context.WithValue(s.ctx, contextKeys.Piid, uuid.FromStringOrNil("843a1ba3-f4b6-4786-9a8c-4914f5d353d3"))
	var product = Product{Name: "name"}

	err := s.repo.UpsertProduct(ctx, &product)
	s.NoError(err)
	s.False(product.Archived)
	s.NotEqual(id, product.ID, "Product is created with different PIID")

	prods, err := s.repo.ListProducts(ctx)
	s.Require().NoError(err)
	s.Len(prods, 1)

	prods, err = s.repo.ListProducts(s.ctx)
	s.Require().NoError(err)
	s.Len(prods, 1)
}
