package product_mgmt

import (
	"context"
	"testing"

	"encore.app/product_mgmt/product"
	"encore.app/product_mgmt/product_instance"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceTest struct {
	suite.Suite
	db      *gorm.DB
	service *Service
	ctx     context.Context
}

func (s *ServiceTest) SetupSuite() {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: prodMgmtDb.Stdlib(),
	}))
	s.Require().NoError(err)
	s.db = db

	s.ctx = s.T().Context()

	app1 := product.App{ID: "app1", Name: "App 1"}
	products := []product.Product{{ID: "product-id", Name: "Product Name", Apps: []product.App{app1}}}
	apps := []product.App{app1}

	pi := product_instance.NewProductInstanceService(db, products, apps)
	prod := product.NewProductService(product.NewProductRepo(products, apps))

	s.service = &Service{pi: pi, prod: prod}
}

func (s *ServiceTest) SetupTest() {
	err := s.db.Exec(`TRUNCATE product_instances;`).Error
	s.Require().NoError(err)
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTest))
}

func (s *ServiceTest) TestCreateInstance() {
	tests := map[string]struct {
		productId   string
		expectError bool
	}{
		"ok":        {productId: "product-id"},
		"not found": {productId: "does not exist", expectError: true},
	}

	for name, test := range tests {
		s.Run(name, func() {
			_, err := s.service.CreateInstance(s.ctx, ProductInstanceParams{
				ProductId:    test.productId,
				InstanceName: "new instance",
			})

			if test.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *ServiceTest) TestFindInstance() {
	resp, err := s.service.CreateInstance(s.ctx, ProductInstanceParams{ProductId: "product-id", InstanceName: "name"})
	s.Require().NoError(err)
	guid := resp.ID
	otherGuid, _ := uuid.NewV4()

	tests := map[string]struct {
		instanceId  uuid.UUID
		expectError bool
	}{
		"ok":        {instanceId: guid},
		"not found": {instanceId: otherGuid, expectError: true},
	}

	for name, test := range tests {
		s.Run(name, func() {
			instance, err := s.service.FindInstance(s.ctx, test.instanceId)

			if test.expectError {
				s.Error(err)
			} else {
				s.Equal(guid, instance.ID)
				s.Equal("product-id", instance.Product.ID)
				s.Equal("Product Name", instance.Product.Name)
				s.Equal("name", instance.Name)
				s.Len(instance.Product.Apps, 1)
				s.Equal("app1", instance.Product.Apps[0].ID)
				s.Equal("App 1", instance.Product.Apps[0].Name)
			}
		})
	}

}

func (s *ServiceTest) TestListInstances() {
	_, err := s.service.CreateInstance(s.ctx, ProductInstanceParams{ProductId: "product-id", InstanceName: "name"})
	s.Require().NoError(err)

	resp, err := s.service.ListInstances(s.ctx)
	s.NoError(err)
	if s.Len(resp.Instances, 1) {
		instance := resp.Instances[0]
		s.Equal("name", instance.Name)
	}
}
