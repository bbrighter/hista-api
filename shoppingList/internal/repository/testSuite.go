package repository

import (
	"context"

	"encore.app/shared/contextKeys"
	"encore.app/shared/generic_queries"
	"encore.app/shoppingList/entity"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ShoppingListTestSuite struct {
	suite.Suite
	ctx         context.Context
	ListRepo    ListRepo
	ItemRepo    ItemRepo
	ProductRepo ProductRepo
}

const GUID_STR = "cf0d4408-8db5-4572-b5d9-4ed873d1341f"

var GUID = uuid.FromStringOrNil(GUID_STR)

func (suite *ShoppingListTestSuite) SetupSuite() {
	suite.ctx = context.WithValue(context.Background(), contextKeys.Piid, GUID)
}

func (suite *ShoppingListTestSuite) SetupSubTest() {
	sqlDb, err := et.NewTestDatabase(suite.ctx, "shopping_list")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}))
	if err != nil {
		panic(err)
	}
	suite.ListRepo = NewListRepo(db)
	suite.ItemRepo = NewItemRepo(db)
	suite.ProductRepo = NewProductRepo(db)
	if err := suite.ItemRepo.db.Debug().Migrator().AutoMigrate(
		&entity.Item{},
		&entity.List{},
		&entity.Product{},
	); err != nil {
		panic(err)
	}
}

func (s *ShoppingListTestSuite) createProduct(ctx context.Context) entity.Product {
	piid, err := generic_queries.PiidFromCtx(ctx)
	s.Require().NoError(err)
	prod := entity.Product{PIID: piid, Name: "name"}
	err = gorm.G[entity.Product](s.ProductRepo.db).Create(ctx, &prod)
	s.Require().NoError(err)
	return prod
}

func (s *ShoppingListTestSuite) createItem() (entity.Product, entity.Item) {
	list := s.createList()
	prod := s.createProduct(s.ctx)
	item := entity.Item{PIID: GUID, ProductId: prod.ID, ProductPiid: GUID, ListId: list.ID, ListPiid: GUID}
	err := gorm.G[entity.Item](s.ItemRepo.db).Create(s.ctx, &item)
	s.Require().NoError(err)
	return prod, item
}

func (s *ShoppingListTestSuite) createList() entity.List {
	piid, err := generic_queries.PiidFromCtx(s.ctx)
	s.Require().NoError(err)
	list := entity.List{PIID: piid}
	err = gorm.G[entity.List](s.ListRepo.db).Create(s.ctx, &list)
	s.Require().NoError(err)
	return list
}

func (s *ShoppingListTestSuite) AssertPostgresError(err error, code string) {
	pgErr, ok := err.(*pgconn.PgError)
	s.True(ok, "expected Postgres error")
	s.Equal(code, pgErr.Code)
}
