package shoppinglist

import "encore.dev/beta/errs"

func (s *ApiTestSuite) TestPatchProductName() {
	s.T().Skip()
	tests := map[string]struct {
		useWrongProductId bool
		useWrongPiid      bool
		name              string
		expectedErrCode   errs.ErrCode
	}{
		"ok":              {name: "new name"},
		"name not unique": {name: "name", expectedErrCode: errs.InvalidArgument},
		"wrong product":   {useWrongProductId: true, expectedErrCode: errs.NotFound},
		"wrong piid":      {useWrongPiid: true, expectedErrCode: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var productId uint = 1000
			if !test.useWrongProductId {
				productId = s.createProduct()
			}
			ctx := s.GetCtx(test.useWrongPiid)

			err := s.service.PatchProductName(ctx, s.piid, productId, PatchProductNameParams{Name: test.name})
			s.assertErrCode(err, test.expectedErrCode)
		})
	}
}

func (s *ApiTestSuite) TestPatchProductArchive() {
	s.T().Skip()
	tests := map[string]struct {
		useWrongProductId bool
		useWrongPiid      bool
		expectedErrCode   errs.ErrCode
	}{
		"ok":            {},
		"wrong product": {useWrongProductId: true, expectedErrCode: errs.NotFound},
		"wrong piid":    {useWrongPiid: true, expectedErrCode: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var productId uint = 1000
			if !test.useWrongProductId {
				productId = s.createProduct()
			}
			ctx := s.GetCtx(test.useWrongPiid)

			err := s.service.PatchArchiveProduct(ctx, s.piid, productId, PatchProductArchiveParams{Archive: true})
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
				productId = s.createProduct()
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
