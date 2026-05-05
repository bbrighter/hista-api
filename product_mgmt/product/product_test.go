package product

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	ctx  context.Context
	prod *ProductRepo
}

func (s *TestSuite) SetupSuite() {
	s.ctx = context.Background()
	newTestProductRepo := func() *ProductRepo {
		app1 := App{ID: "app1", Name: "App 1"}
		app2 := App{ID: "app2", Name: "App 2"}
		return &ProductRepo{
			apps: []App{app1, app2},
			products: []Product{
				{ID: "test-id", Name: "test-name", Apps: []App{app1}},
			},
		}
	}
	s.prod = newTestProductRepo()
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestProductList() {
	products := s.prod.ListProducts()

	s.Len(products, 1)

	var histaProduct Product
	for _, p := range products {
		if p.ID == "test-id" {
			histaProduct = p
			break
		}
	}
	s.Equal("test-name", histaProduct.Name)
	s.Len(histaProduct.Apps, 1)

	var userMgmtApp App
	for _, a := range histaProduct.Apps {
		if a.ID == "app1" {
			userMgmtApp = a
			break
		}
	}
	s.Equal("App 1", userMgmtApp.Name)
}
