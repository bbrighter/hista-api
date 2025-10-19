package repository

import (
	"context"

	"encore.app/shared/contextKeys"
	"encore.app/shared/generic_queries"
	"encore.app/shoppingList/entity"
	"encore.dev/types/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/suite"
	postgresContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresContainer struct {
	*postgresContainer.PostgresContainer
	ConnectionString string
}

func NewPostgresContainer() *PostgresContainer {
	ctx := context.Background()
	pgContainer, err := postgresContainer.Run(ctx,
		"postgres:15.3-alpine",
		postgresContainer.BasicWaitStrategies(),
	)
	if err != nil {
		panic(err)
	}

	connString, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connString,
	}
}

type ShoppingListTestSuite struct {
	suite.Suite
	pgContainer *PostgresContainer
	ctx         context.Context
	ListRepo    ListRepo
	ItemRepo    ItemRepo
	ProductRepo ProductRepo
}

const GUID_STR = "cf0d4408-8db5-4572-b5d9-4ed873d1341f"

var GUID = uuid.FromStringOrNil(GUID_STR)

func (suite *ShoppingListTestSuite) SetupSuite() {
	suite.ctx = context.WithValue(context.Background(), contextKeys.Piid, GUID)

	suite.pgContainer = NewPostgresContainer()

	db, err := gorm.Open(postgres.Open(suite.pgContainer.ConnectionString))
	if err != nil {
		panic(err)
	}

	suite.ListRepo = NewListRepo(db)
	suite.ItemRepo = NewItemRepo(db)
	suite.ProductRepo = NewProductRepo(db)
}

func (suite *ShoppingListTestSuite) SetupSubTest() {
	if err := suite.ItemRepo.db.Migrator().AutoMigrate(
		&entity.Item{},
		&entity.List{},
		&entity.Product{},
	); err != nil {
		panic(err)
	}
}

func (suite *ShoppingListTestSuite) TearDownSubTest() {
	if err := suite.ItemRepo.db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error; err != nil {
		panic(err)
	}
}

func (suite *ShoppingListTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
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
