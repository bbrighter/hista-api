package repository

import (
	"context"

	"encore.app/shared/contextKeys"
	"encore.app/shoppingList/entity"
	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

func (s *RepoTestSuite) TestFirstList() {
	tests := map[string]struct {
		useWrongId   bool
		useWrongPiid bool
		deleteBefore bool
		expectError  error
	}{
		"ok":              {},
		"wrong piid":      {useWrongPiid: true, expectError: gorm.ErrRecordNotFound},
		"already deleted": {deleteBefore: true, expectError: gorm.ErrRecordNotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			initList, _, _ := s.createItem()
			listId := initList.ID
			if test.useWrongId {
				listId = 1000
			}
			ctx := s.ctx
			if test.useWrongPiid {
				piid, _ := uuid.NewV4()
				ctx = context.WithValue(ctx, contextKeys.Piid, piid)
			}
			if test.deleteBefore {
				_, err := gorm.G[entity.List](s.db).Where("id = ?", listId).Delete(s.ctx)
				s.Require().NoError(err)
			}

			list, err := s.ListRepo.First(ctx)
			if test.expectError != nil {
				s.ErrorIs(err, test.expectError)
				return
			}
			s.Greater(list.ID, uint(0))
			s.Len(list.Items, 1)
		})
	}
}
