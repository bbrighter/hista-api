package shoppinglist

import "encore.dev/beta/errs"

func (s *ApiTestSuite) TestDeleteList() {
	tests := map[string]struct {
		useWrongId        bool
		statusCode        errs.ErrCode
		hasUncheckedItems bool
	}{
		"ok":             {},
		"not found":      {useWrongId: true, statusCode: errs.NotFound},
		"has items left": {hasUncheckedItems: true, statusCode: errs.InvalidArgument},
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
			err := s.service.DeleteList(s.ctx, GUID, listId)
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
