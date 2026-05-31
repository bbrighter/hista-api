package shoppinglist

import (
	"encore.dev/beta/errs"
	"encore.dev/types/option"
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
				productId = s.createProduct("name")
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

func (s *ApiTestSuite) TestPutItemByName() {
	tests := map[string]struct {
		nameAlreadyExists bool
		useWrongListId    bool
		expectedErrorCode errs.ErrCode
		useWrongPiid      bool
	}{
		"ok":              {},
		"list not found":  {useWrongListId: true, expectedErrorCode: errs.NotFound},
		"name not exists": {nameAlreadyExists: true},
		"use wrong piid":  {useWrongPiid: true, expectedErrorCode: errs.NotFound},
	}

	for name, test := range tests {
		s.Run(name, func() {
			if test.nameAlreadyExists {
				s.createProduct("name")
			}
			var listId uint = 1000
			if !test.useWrongListId {
				listId = s.createList()
			}

			ctx := s.GetCtx(test.useWrongPiid)

			resp, err := s.service.PutItemByName(ctx, s.piid, listId, ItemNameParams{Name: "name"})
			s.assertErrCode(err, test.expectedErrorCode)
			if test.expectedErrorCode == 0 {
				s.Greater(resp.ID, uint(0))
			}
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
				itemId = s.createItem(listId).ID
			}
			ctx := s.GetCtx(false)
			err := s.service.DeleteItem(ctx, s.piid, itemId)

			s.assertErrCode(err, test.expectedErrCode)
		})
	}
}

func (s *ApiTestSuite) TestPatchItem() {
	tests := map[string]struct {
		useWrongItemId  bool
		quantity        option.Option[uint8]
		checked         option.Option[bool]
		expectedErrCode errs.ErrCode
	}{
		"10":        {quantity: option.Some(uint8(10))},
		"0 as nil":  {quantity: option.Some(uint8(0))},
		"nil":       {},
		"not found": {useWrongItemId: true, expectedErrCode: errs.NotFound},
		"checked":   {checked: option.Some(true)},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var itemId uint = 1000
			if !test.useWrongItemId {
				listId := s.createList()
				itemId = s.createItem(listId).ID
			}
			ctx := s.GetCtx(false)
			params := ItemPatchParams{Quantity: test.quantity}
			err := s.service.PatchItem(ctx, s.piid, itemId, params)
			s.assertErrCode(err, test.expectedErrCode)

			var respItem ItemResponse
			if test.quantity.IsSome() {
				list, err := s.service.PostOrGetList(ctx, s.piid)
				s.Require().NoError(err)

				for _, item := range list.Items {
					if item.ID == itemId {
						respItem = item
						break
					}
				}
				if test.quantity.GetOrElse(0) == 0 {
					s.Equal(option.None[uint8](), respItem.Quantity)
				}
				if test.quantity.GetOrElse(0) > 0 {
					s.Equal(test.quantity, respItem.Quantity)
				}
				if test.checked.IsSome() {
					s.Equal(test.checked.MustGet(), respItem.Checked)
				}
			}
		})
	}
}
