package internal

import (
	"encore.app/shoppingList/entity"
	"gorm.io/gorm"
)

func (s *internalTestSuite) TestDeleteList() {
	ctx := s.ctx
	tests := map[string]struct {
		items           []entity.Item
		listItemError   error
		deleteListError error
		expectedError   string
	}{
		"ok":              {items: []entity.Item{{ID: 1, Checked: true}}},
		"no items":        {items: []entity.Item{}},
		"unchecked items": {items: []entity.Item{{ID: 1, Checked: false}, {ID: 2, Checked: true}}, expectedError: "item unchecked"},
		"delete error":    {deleteListError: gorm.ErrRecordNotFound, expectedError: "not found"},
		"list item error": {listItemError: gorm.ErrRecordNotFound, expectedError: "not found"},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var listId uint = 1
			s.itemRepo.On("List", ctx, listId).Return(test.items, test.listItemError)
			s.listRepo.On("Delete", ctx, listId).Return(test.deleteListError)

			err := s.listUc.Delete(ctx, listId)
			if test.expectedError != "" {
				s.Error(err)
				s.ErrorContains(err, test.expectedError)
			} else {
				s.NoError(err)
				s.listRepo.AssertCalled(s.T(), "Delete", ctx, listId)
			}
		})
	}
}

func (s *internalTestSuite) TestCreateOrGetList() {
	ctx := s.ctx
	tests := map[string]struct {
		listExistsError    error
		expectCreateCalled bool
		expectedError      error
	}{
		"ok, list exists":        {},
		"ok, list doesn't exits": {listExistsError: gorm.ErrRecordNotFound, expectCreateCalled: true},
	}

	for name, test := range tests {
		s.Run(name, func() {

			s.listRepo.On("First", ctx).Return(&entity.List{ID: 1}, test.listExistsError)
			s.listRepo.On("Create", ctx).Return(2, nil)

			list, err := s.listUc.FirstOrCreate(ctx)
			if test.expectedError != nil {
				s.Error(err)
				return
			}
			s.NoError(err)

			if test.expectCreateCalled {
				s.listRepo.AssertCalled(s.T(), "Create", ctx)
				s.EqualValues(list.ID, 2)
			} else {
				s.EqualValues(list.ID, 1)
			}
		})
	}
}
