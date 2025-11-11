package shoppinglist

import (
	"context"
	"testing"

	"encore.app/shared/contextKeys"
	entity "encore.app/shoppingList/entity"
	"encore.dev/beta/errs"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ApiTestSuite struct {
	suite.Suite
	service *Service
	ctx     context.Context
	piid    uuid.UUID
	db      *gorm.DB
}

func (suite *ApiTestSuite) SetupSuite() {
	guid, err := uuid.FromString("2012b8a8-df7f-407c-bda1-9567b5b8f06d")
	suite.Require().NoError(err)

	suite.piid = guid
	suite.ctx = context.WithValue(context.Background(), contextKeys.Piid, guid)
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
	resp, err := suite.service.GetOrCreateList(suite.ctx, suite.piid)
	suite.Require().NoError(err)
	return resp.ID
}

func (suite *ApiTestSuite) createProduct() uint {
	var product = entity.Product{Name: "name", PIID: suite.piid}
	err := gorm.G[entity.Product](suite.db).Create(suite.ctx, &product)
	suite.Require().NoError(err)
	return product.ID
}

func (suite *ApiTestSuite) createItem(listId uint) uint {
	productId := suite.createProduct()
	var item = entity.Item{PIID: suite.piid, ProductPiid: suite.piid, ListPiid: suite.piid, ProductId: productId, ListId: listId}
	err := gorm.G[entity.Item](suite.db).Create(suite.ctx, &item)
	suite.Require().NoError(err)
	return item.ID
}

func (suite *ApiTestSuite) assertErrCode(err error, expectedCode errs.ErrCode) {
	if expectedCode == 0 {
		return
	}
	suite.Require().NotNil(err)
	if err != nil {
		print(err.Error())
	}
	encoreErr, ok := err.(*errs.Error)
	suite.Require().True(ok)
	suite.Equal(encoreErr.Code, expectedCode)
}

func (suite *ApiTestSuite) GetCtx(different bool) context.Context {
	if different {
		piid, err := uuid.FromString("d771f40a-b291-4671-a67c-edcb0181ebc7")
		suite.Require().NoError(err)
		return context.WithValue(suite.ctx, contextKeys.Piid, piid)
	}
	return suite.ctx

}
