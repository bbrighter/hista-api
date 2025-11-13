package shoppinglist

import (
	entity "encore.app/shoppingList/entity"
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

func (s *ApiTestSuite) TestDeleteItem() {
	tests := map[string]struct {
		useWrongItemId  bool
		expectedErrCode errs.ErrCode
	}{
		"ok":        {},
		"not found": {useWrongItemId: true, expectedErrCode: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var itemId uint = 1000
			if !test.useWrongItemId {
				listId := s.createList()
				itemId = s.createItem(listId)
			}
			ctx := s.GetCtx(false)
			err := s.service.DeleteItem(ctx, s.piid, itemId)

			s.assertErrCode(err, test.expectedErrCode)
		})
	}
}

func (s *ApiTestSuite) TestPatchItem() {
	var quantity0, quantity10 uint8
	quantity0 = 0
	quantity10 = 10
	tests := map[string]struct {
		useWrongItemId  bool
		quantity        *uint8
		expectedErrCode errs.ErrCode
	}{
		"10":        {quantity: &quantity10},
		"0 as nil":  {quantity: &quantity0},
		"nil":       {quantity: nil},
		"not found": {useWrongItemId: true, expectedErrCode: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var itemId uint = 1000
			if !test.useWrongItemId {
				listId := s.createList()
				itemId = s.createItem(listId)
			}
			ctx := s.GetCtx(false)
			params := ItemPatchParams{Quantity: test.quantity}
			err := s.service.PatchItem(ctx, s.piid, itemId, params)
			s.assertErrCode(err, test.expectedErrCode)

			var dbItem entity.ItemResponse
			if test.quantity != nil {
				resp, _ := s.service.GetOrCreateList(ctx, s.piid)
				for _, item := range resp.Items {
					if item.ID == itemId {
						dbItem = item
						break
					}
				}
				if *test.quantity == 0 {
					s.Nil(dbItem.Quantity)
				}
				if *test.quantity > 0 {
					s.Equal(test.quantity, dbItem.Quantity)
				}
			}
		})
	}
}
