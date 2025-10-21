package shoppinglist

import (
	"context"
	"testing"

	"encore.app/shared/contextKeys"
	entity "encore.app/shoppingList/entity"
	"encore.dev/et"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ApiTestSuite struct {
	suite.Suite
	service *Service
	ctx     context.Context
	db      *gorm.DB
}

func (suite *ApiTestSuite) SetupSuite() {
	suite.ctx = context.WithValue(context.Background(), contextKeys.Piid, GUID)
	sqlDb, err := et.NewTestDatabase(suite.ctx, "shopping_list")
	suite.Require().NoError(err)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}))
	suite.db = db.Debug()
	suite.Require().NoError(err)
	suite.service = initServiceWithDb(suite.db)
}

func (suite *ApiTestSuite) TearDownSubTest() {
	var err error
	tx := suite.db
	err = tx.Exec(`DELETE FROM items`).Error
	suite.Require().NoError(err)
	err = tx.Exec(`DELETE FROM lists`).Error
	suite.Require().NoError(err)
	err = tx.Exec(`DELETE FROM products`).Error
	suite.Require().NoError(err)
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func (suite *ApiTestSuite) createList() uint {
	resp, err := suite.service.GetOrCreateList(suite.ctx, GUID)
	suite.Require().NoError(err)
	return resp.ID
}

func (suite *ApiTestSuite) createProduct() uint {
	var product = entity.Product{Name: "name", PIID: GUID}
	err := gorm.G[entity.Product](suite.db).Create(suite.ctx, &product)
	suite.Require().NoError(err)
	return product.ID
	// resp, err := suite.service.PostItemByName(suite.ctx, ItemNameParams{Name: "name"})
	// suite.Require().NoError(err)
}

func (suite *ApiTestSuite) createItem(listId uint) uint {
	productId := suite.createProduct()
	var item = entity.Item{PIID: GUID, ProductPiid: GUID, ListPiid: GUID, ProductId: productId, ListId: listId}
	err := gorm.G[entity.Item](suite.db).Create(suite.ctx, &item)
	suite.Require().NoError(err)
	return item.ID
}
