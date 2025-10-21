package repository

import (
	"context"
	"testing"

	"encore.app/shared/contextKeys"
	"encore.app/shared/generic_queries"
	"encore.app/shoppingList/entity"
	"encore.app/shoppingList/internal"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RepoTestSuite struct {
	suite.Suite
	ctx         context.Context
	ListRepo    internal.ListRepo
	ItemRepo    internal.ItemRepo
	ProductRepo internal.ProductRepo
	db          *gorm.DB
}

const GUID_STR = "cf0d4408-8db5-4572-b5d9-4ed873d1341f"

var GUID = uuid.FromStringOrNil(GUID_STR)

func (suite *RepoTestSuite) SetupSuite() {
	suite.ctx = context.WithValue(context.Background(), contextKeys.Piid, GUID)
	sqlDb, err := et.NewTestDatabase(suite.ctx, "shopping_list")
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
	suite.ListRepo = NewListRepo(db)
	suite.ItemRepo = NewItemRepo(db)
	suite.ProductRepo = NewProductRepo(db)
}

func (suite *RepoTestSuite) SetupSubTest() {
	if err := suite.db.Migrator().AutoMigrate(
		&entity.Item{},
		&entity.List{},
		&entity.Product{},
	); err != nil {
		panic(err)
	}
}

func (suite *RepoTestSuite) TearDownSubTest() {
	var err error
	tx := suite.db
	err = tx.Exec(`DELETE FROM items`).Error
	suite.Require().NoError(err)
	err = tx.Exec(`DELETE FROM lists`).Error
	suite.Require().NoError(err)
	err = tx.Exec(`DELETE FROM products`).Error
	suite.Require().NoError(err)
}

func (s *RepoTestSuite) createProduct(ctx context.Context) entity.Product {
	piid, err := generic_queries.PiidFromCtx(ctx)
	s.Require().NoError(err)
	prod := entity.Product{PIID: piid, Name: "name"}
	err = gorm.G[entity.Product](s.db).Create(ctx, &prod)
	s.Require().NoError(err)
	return prod
}

func (s *RepoTestSuite) createItem() (entity.List, entity.Product, entity.Item) {
	list := s.createList()
	prod := s.createProduct(s.ctx)
	item := entity.Item{PIID: GUID, ProductId: prod.ID, ProductPiid: GUID, ListId: list.ID, ListPiid: GUID}
	err := gorm.G[entity.Item](s.db).Create(s.ctx, &item)
	s.Require().NoError(err)
	return list, prod, item
}

func (s *RepoTestSuite) createList() entity.List {
	piid, err := generic_queries.PiidFromCtx(s.ctx)
	s.Require().NoError(err)
	list := entity.List{PIID: piid}
	err = gorm.G[entity.List](s.db).Create(s.ctx, &list)
	s.Require().NoError(err)
	return list
}

func (s *RepoTestSuite) AssertPostgresError(err error, code string) {
	pgErr, ok := err.(*pgconn.PgError)
	s.True(ok, "expected Postgres error")
	s.Equal(code, pgErr.Code)
}

func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}
