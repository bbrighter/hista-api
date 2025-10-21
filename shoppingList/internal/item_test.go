package internal

import (
	"encore.app/shoppingList/entity"
	"gorm.io/gorm"
)

func (s *internalTestSuite) TestAddItemByName() {
	ctx := s.ctx
	tests := map[string]struct {
		prodName      string
		prodCreateErr error
	}{
		"ok":       {prodName: "name"},
		"trimming": {prodName: "  name  "},
		"prod err": {prodName: "name", prodCreateErr: gorm.ErrInvalidDB},
	}

	for name, test := range tests {
		s.Run(name, func() {
			s.prodRepo.On("Create", ctx, "name").Return(1, test.prodCreateErr)
			s.itemRepo.On("Create", ctx, uint(1), uint(1)).Return(10, nil)
			s.itemRepo.On("Find", ctx, uint(10)).Return(&entity.Item{ID: 10, ProductId: 1}, nil)

			item, err := s.itemUc.AddItemByName(ctx, 1, test.prodName)
			if test.prodCreateErr != nil {
				s.itemRepo.AssertNotCalled(s.T(), "Create")
				return
			}
			s.prodRepo.AssertCalled(s.T(), "Create", ctx, "name")
			s.NoError(err)
			s.EqualValues(10, item.ID)
		})
	}
}
