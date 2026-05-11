package shoppinglist

import (
	"encore.dev/beta/errs"
)

func (s *ApiTestSuite) TestDeleteList() {
	tests := map[string]struct {
		useWrongId        bool
		useWrongPiid      bool
		statusCode        errs.ErrCode
		hasUncheckedItems bool
		forceDelete       bool
	}{
		"ok":             {},
		"not found":      {useWrongId: true, statusCode: errs.NotFound},
		"has items left": {hasUncheckedItems: true, statusCode: errs.InvalidArgument},
		"wrong piid":     {useWrongPiid: true, statusCode: errs.NotFound},
		"force delete":   {hasUncheckedItems: true, forceDelete: true},
	}
	for name, test := range tests {
		s.Run(name, func() {
			listId := s.createList()
			if test.useWrongId {
				listId = 1000
			}
			if test.hasUncheckedItems {
				s.createItem(listId)
			}
			ctx := s.GetCtx(test.useWrongPiid)
			err := s.service.DeleteList(ctx, s.piid, listId, DeleteListForceDeleteParam{Force: test.forceDelete})
			if test.statusCode != 0 {
				encoreErr, ok := err.(*errs.Error)
				s.True(ok)
				s.Equal(test.statusCode, encoreErr.Code)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *ApiTestSuite) TestCreateList() {
	tests := map[string]struct {
		listExists        bool
		useWrongPiid      bool
		newIdIsOldId      bool
		expectedErrorCode errs.ErrCode
	}{
		"no list exists": {listExists: false},
		"list exists":    {listExists: true, newIdIsOldId: true},
		"wrong piid":     {useWrongPiid: true},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var listId uint = 0
			if test.listExists {
				listId = s.createList()
			}
			ctx := s.GetCtx(test.useWrongPiid)
			list, err := s.service.PostOrGetList(ctx, s.piid)
			if test.expectedErrorCode > 0 {
				s.assertErrCode(err, test.expectedErrorCode)
				return
			}
			s.NoError(err)
			if test.newIdIsOldId {
				s.Equal(listId, list.ID)
			} else {
				s.Greater(list.ID, uint(0))
			}
		})
	}
}
