package shoppinglist

import (
	"context"
	"fmt"
	"testing"

	"encore.app/shared/contextKeys"
	shoppinglist "encore.app/shoppingList/internal/shoppingList"
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
}

func (suite *ApiTestSuite) TearDownSubTest() {
	tables := []string{"items", "lists", "products"}
	for _, table := range tables {
		err := suite.db.Exec(fmt.Sprintf(`DELETE FROM "%s"`, table)).Error
		suite.Require().NoError(err)
	}
}

func (s *ApiTestSuite) SetupTest() {
	s.service = initServiceWithDb(s.db)
}

func (suite *ApiTestSuite) SetupSubTest() {
	suite.service = initServiceWithDb(suite.db)
	// suite.createMoment(time.Date(2020, 5, 3, 2, 1, 0, 0, time.UTC))
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func (suite *ApiTestSuite) createList() uint {
	resp, err := suite.service.PostOrGetList(suite.ctx, suite.piid)
	suite.Require().NoError(err)
	return resp.ID
}

func (suite *ApiTestSuite) createProduct(name string) uint {
	var product = shoppinglist.Product{Name: name, PIID: suite.piid}
	err := gorm.G[shoppinglist.Product](suite.db).Create(suite.ctx, &product)
	suite.Require().NoError(err)
	return product.ID
}

func (suite *ApiTestSuite) createItem(listId uint) shoppinglist.Item {
	productId := suite.createProduct("name")
	var item = shoppinglist.Item{PIID: suite.piid, ProductPiid: suite.piid, ListPiid: suite.piid, ProductId: productId, ListId: listId}
	err := gorm.G[shoppinglist.Item](suite.db).Create(suite.ctx, &item)
	suite.Require().NoError(err)
	return item
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

func (s *ApiTestSuite) TestShoppingListWorkflow() {
	// Moments are 0
	moments, err := s.service.GetMoments(s.ctx, s.piid)
	s.NoError(err)
	s.Equal(0, moments.Items)
	s.Equal(0, moments.Products)

	// Create list
	list, err := s.service.PostOrGetList(s.ctx, s.piid)
	s.NoError(err)
	s.Len(list.Items, 0)
	listId := list.ID
	s.NotEqual(0, listId)

	// Items moments are +1
	moments, err = s.service.GetMoments(s.ctx, s.piid)
	s.NoError(err)
	s.Equal(1, moments.Items)
	s.Equal(0, moments.Products)

	// Add products and item
	itemResp, err := s.service.PostItemByName(s.ctx, s.piid, listId, ItemNameParams{Name: "product"})
	s.NoError(err)
	productId := itemResp.ProductId
	s.NotEqual(0, productId)
	itemId := itemResp.ID
	s.NotEqual(0, itemId)

	// Both moments are +1
	moments, err = s.service.GetMoments(s.ctx, s.piid)
	s.NoError(err)
	s.Equal(2, moments.Items)
	s.Equal(1, moments.Products)

	// Increase item quantity
	var newQuantity uint8 = 3
	err = s.service.PatchItem(s.ctx, s.piid, itemId, ItemPatchParams{Quantity: option.Some(newQuantity)})
	s.NoError(err)

	// Moments for items are +1
	moments, err = s.service.GetMoments(s.ctx, s.piid)
	s.NoError(err)
	s.Equal(3, moments.Items)
	s.Equal(1, moments.Products)

	// Check item
	err = s.service.PatchItem(s.ctx, s.piid, itemId, ItemPatchParams{Checked: option.Some(true)})
	s.NoError(err)

	// Moments for items are +1
	moments, err = s.service.GetMoments(s.ctx, s.piid)
	s.NoError(err)
	s.Equal(4, moments.Items)
	s.Equal(1, moments.Products)

	// Delete list
	err = s.service.DeleteList(s.ctx, s.piid, listId, DeleteListForceDeleteParam{Force: true})
	s.NoError(err)

	// Moments for items are +1
	moments, err = s.service.GetMoments(s.ctx, s.piid)
	s.NoError(err)
	s.Equal(5, moments.Items)
	s.Equal(1, moments.Products)
}
