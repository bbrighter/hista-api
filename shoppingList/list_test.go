package shoppinglist

import (
	"encore.dev/beta/errs"
)

func (s *ApiTestSuite) TestGetOrCreateList() {
	tests := map[string]struct {
		listExists   bool
		useWrongPiid bool
		newIdIsOldId bool
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
			list, err := s.service.GetOrCreateList(ctx, s.piid)
			s.NoError(err)
			if test.newIdIsOldId {
				s.Equal(listId, list.ID)
			} else {
				s.Greater(list.ID, uint(0))
			}
		})
	}
}

func (s *ApiTestSuite) TestDeleteList() {
	tests := map[string]struct {
		useWrongId        bool
		useWrongPiid      bool
		statusCode        errs.ErrCode
		hasUncheckedItems bool
	}{
		"ok":             {},
		"not found":      {useWrongId: true, statusCode: errs.NotFound},
		"has items left": {hasUncheckedItems: true, statusCode: errs.InvalidArgument},
		"wrong piid":     {useWrongPiid: true, statusCode: errs.NotFound},
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
			err := s.service.DeleteList(ctx, s.piid, listId)
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
