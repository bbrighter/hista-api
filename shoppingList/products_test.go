package shoppinglist

import (
	"encore.dev/beta/errs"
	"encore.dev/types/option"
)

func (s *ApiTestSuite) TestPatchProduct() {
	tests := map[string]struct {
		useWrongProductId bool
		useWrongPiid      bool
		name              option.Option[string]
		archived          option.Option[bool]
		expectedErrCode   errs.ErrCode
	}{
		"ok":              {name: option.Some("new name")},
		"name not unique": {name: option.Some("existing name"), expectedErrCode: errs.AlreadyExists},
		"archived":        {archived: option.Some(true)},
		"wrong product":   {useWrongProductId: true, expectedErrCode: errs.NotFound},
		"wrong piid":      {useWrongPiid: true, expectedErrCode: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			s.createProduct("existing name")
			var productId uint = 1000
			if !test.useWrongProductId {
				productId = s.createProduct("name")
			}
			ctx := s.GetCtx(test.useWrongPiid)

			err := s.service.PatchProduct(ctx, s.piid, productId, PatchProductParams{Name: test.name})
			s.assertErrCode(err, test.expectedErrCode)
		})
	}
}

func (s *ApiTestSuite) TestDeleteProduct() {
	s.T().Skip()
	tests := map[string]struct {
		isUsedInItem      bool
		useWrongProductId bool
		useWrongPiid      bool
		expectedErrCode   errs.ErrCode
	}{
		"ok":            {},
		"in use":        {isUsedInItem: true, expectedErrCode: errs.AlreadyExists},
		"wrong product": {useWrongProductId: true, expectedErrCode: errs.NotFound},
		"wrong piid":    {useWrongPiid: true, expectedErrCode: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var productId uint = 1000
			if !test.useWrongProductId {
				productId = s.createProduct("name")
			}
			if test.isUsedInItem {
				listId := s.createList()
				productId = s.createItem(listId).ProductId
			}
			ctx := s.GetCtx(test.useWrongPiid)

			err := s.service.DeleteProduct(ctx, s.piid, productId)
			s.assertErrCode(err, test.expectedErrCode)
		})
	}
}
