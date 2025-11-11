package repository

import (
	"context"

	"encore.app/shared/contextKeys"
	"encore.app/shared/generic_queries"
	"encore.app/shoppingList/entity"
	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

func (s *RepoTestSuite) TestCreateItem() {
	tests := map[string]struct {
		useNonExistingId    bool
		useWrongPiid        bool
		createTwice         bool
		expectPostgresError string
	}{
		"ok":            {},
		"wrong prod id": {useNonExistingId: true, expectPostgresError: "23503"},
		"wrong piid":    {useWrongPiid: true, expectPostgresError: "23503"},
		"create twice":  {createTwice: true, expectPostgresError: "23505"},
	}

	for name, test := range tests {
		s.Run(name, func() {
			var prod entity.Product = s.createProduct(s.ctx)
			var list entity.List = s.createList()
			prodId := prod.ID
			if test.useNonExistingId {
				prodId = 1000
			}
			createCtx := s.ctx
			if test.useWrongPiid {
				newGuid, _ := uuid.NewV4()
				createCtx = context.WithValue(s.ctx, contextKeys.Piid, newGuid)
			}
			if test.createTwice {
				s.ItemRepo.Create(createCtx, prodId, list.ID)
			}
			itemId, err := s.ItemRepo.Create(createCtx, prodId, list.ID)
			if test.expectPostgresError != "" {
				s.AssertPostgresError(err, test.expectPostgresError)
			} else {
				s.Greater(itemId, uint(0))
			}
		})
	}
}

func (s *RepoTestSuite) TestDeleteItem() {
	tests := map[string]struct {
		useNonExistingId bool
		useWrongPiid     bool
		expectError      error
	}{
		"ok":             {},
		"id not found":   {useNonExistingId: true, expectError: gorm.ErrRecordNotFound},
		"piid not found": {useWrongPiid: true, expectError: gorm.ErrRecordNotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			_, _, item := s.createItem()
			var id []uint = []uint{item.ID}
			if test.useNonExistingId {
				id = []uint{1000}
			}
			ctx := s.ctx
			if test.useWrongPiid {
				guid, _ := uuid.NewV4()
				ctx = context.WithValue(ctx, contextKeys.Piid, guid)
			}

			err := s.ItemRepo.Delete(ctx, id)
			if test.expectError != nil {
				s.ErrorIs(err, test.expectError)
				return
			}
			s.NoError(err)
			var count int64
			s.db.Unscoped().Debug().Model(&entity.Item{}).Count(&count)
			s.EqualValues(count, 0)
		})
	}
}

func (s *RepoTestSuite) TestCheckItem() {
	tests := map[string]struct {
		numberOfChecks   int
		useNonExistingId bool
		useWrongPiid     bool
		expectedValue    bool
		expectError      error
	}{
		"ok":          {numberOfChecks: 1, expectedValue: true},
		"check twice": {numberOfChecks: 2, expectedValue: false},
		"not found":   {numberOfChecks: 1, useNonExistingId: true, expectError: gorm.ErrRecordNotFound},
		"wrong piid":  {numberOfChecks: 1, useWrongPiid: true, expectError: gorm.ErrRecordNotFound},
	}

	for name, test := range tests {
		s.Run(name, func() {
			_, _, item := s.createItem()

			itemId := item.ID
			if test.useNonExistingId {
				itemId = 1000
			}
			ctx := s.ctx
			if test.useWrongPiid {
				newGuid, _ := uuid.NewV4()
				ctx = context.WithValue(s.ctx, contextKeys.Piid, newGuid)
			}

			for range test.numberOfChecks {
				err := s.ItemRepo.Check(ctx, itemId)
				if test.expectError != nil {
					s.Error(err)
					return
				}
				s.NoError(err)
			}
			respItem, err := generic_queries.First[*entity.Item](ctx, s.db, itemId)
			s.Require().NoError(err)
			s.Equal(test.expectedValue, respItem.Checked)
		})
	}
}

func (s *RepoTestSuite) TestListItems() {
	tests := map[string]struct {
		useWrongPiid bool
		useWrongList bool
		expectedLen  int
	}{
		"ok":         {expectedLen: 1},
		"wrong piid": {useWrongPiid: true, expectedLen: 0},
		"wrong list": {useWrongList: true, expectedLen: 0},
	}

	for name, test := range tests {
		s.Run(name, func() {
			list, _, _ := s.createItem()

			ctx := s.ctx
			if test.useWrongPiid {
				guid, _ := uuid.NewV4()
				ctx = context.WithValue(ctx, contextKeys.Piid, guid)
			}
			listId := list.ID
			if test.useWrongList {
				listId = 1000
			}

			items, err := s.ItemRepo.List(ctx, listId)
			s.NoError(err)
			s.Len(items, test.expectedLen)
		})
	}
}
