package shoppinglist

import (
	"encore.dev/beta/errs"
)

func (s *ApiTestSuite) TestPostItem() {
	tests := map[string]struct {
		useWrongProductId bool
		useWrongListId    bool
		expectedErrorCode errs.ErrCode
		useWrongPiid      bool
	}{
		"ok":         {},
		"no product": {useWrongProductId: true, expectedErrorCode: errs.NotFound},
		"not list":   {useWrongListId: true, expectedErrorCode: errs.NotFound},
		"wrong piid": {useWrongPiid: true, expectedErrorCode: errs.NotFound},
	}

	for name, test := range tests {
		s.Run(name, func() {
			var productId uint = 1000
			if !test.useWrongProductId {
				productId = s.createProduct()
			}
			var listId uint = 1000
			if !test.useWrongListId {
				listId = s.createList()
			}
			ctx := s.GetCtx(test.useWrongPiid)

			resp, err := s.service.PostItem(ctx, s.piid, listId, productId)
			s.assertErrCode(err, test.expectedErrorCode)
			if test.expectedErrorCode == 0 {
				s.Greater(resp.ID, uint(0))
			}
		})
	}
}

func (s *ApiTestSuite) TestPostItemByName() {
	tests := map[string]struct {
		nameAlreadyExists bool
		useWrongListId    bool
		expectedErrorCode errs.ErrCode
		useWrongPiid      bool
	}{
		"ok":              {},
		"list not found":  {useWrongListId: true, expectedErrorCode: errs.NotFound},
		"name not exists": {nameAlreadyExists: true, expectedErrorCode: errs.AlreadyExists},
		"use wrong piid":  {useWrongPiid: true, expectedErrorCode: errs.NotFound},
	}

	for name, test := range tests {
		s.Run(name, func() {
			if test.nameAlreadyExists {
				s.createProduct()
			}
			var listId uint = 1000
			if !test.useWrongListId {
				listId = s.createList()
			}

			ctx := s.GetCtx(test.useWrongPiid)

			resp, err := s.service.PostItemByName(ctx, s.piid, listId, ItemNameParams{Name: "name"})
			s.assertErrCode(err, test.expectedErrorCode)
			if test.expectedErrorCode == 0 {
				s.Greater(resp.ID, uint(0))
			}
		})
	}
}

func (s *ApiTestSuite) TestCheckItem() {
	tests := map[string]struct {
		useWrongItemId  bool
		expectedErrCode errs.ErrCode
		useWrongPiid    bool
	}{
		"ok":         {},
		"not found":  {useWrongItemId: true, expectedErrCode: errs.NotFound},
		"wrong piid": {useWrongPiid: true, expectedErrCode: errs.NotFound},
	}

	for name, test := range tests {
		s.Run(name, func() {
			var itemId uint = 1000
			if !test.useWrongItemId {
				listId := s.createList()
				itemId = s.createItem(listId)
			}
			ctx := s.GetCtx(test.useWrongPiid)
			err := s.service.CheckItem(ctx, s.piid, itemId)

			s.assertErrCode(err, test.expectedErrCode)
		})
	}
}
