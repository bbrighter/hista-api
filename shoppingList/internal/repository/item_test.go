package repository

import (
	"context"

	"encore.app/shared/contextKeys"
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
			s.db.Unscoped().Model(&entity.Item{}).Count(&count)
			s.EqualValues(count, 0)
		})
	}
}

func (s *RepoTestSuite) TestPatchItem() {
	var defaultQuantity uint8 = 3
	defaultMap := map[string]any{"quantity": &defaultQuantity}
	var nilQuantity *uint8
	tests := map[string]struct {
		useNonExistingId bool
		useWrongPiid     bool
		expectedError    error
		patchMap         map[string]any
	}{
		"ok":             {patchMap: defaultMap},
		"not found":      {useWrongPiid: true, patchMap: defaultMap, expectedError: gorm.ErrRecordNotFound},
		"wrong piid":     {useWrongPiid: true, patchMap: defaultMap, expectedError: gorm.ErrRecordNotFound},
		"empty quantity": {patchMap: map[string]any{"quantity": nilQuantity}},
		"check":          {patchMap: map[string]any{"checked": gorm.Expr("NOT checked")}},
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

			err := s.ItemRepo.Patch(ctx, itemId, test.patchMap)
			if test.expectedError != nil {
				s.Error(err, test.expectedError)
				return
			}
			s.NoError(err)
			dbItem, err := s.ItemRepo.Find(ctx, itemId)
			s.Require().NoError(err)
			if _, ok := test.patchMap["quantity"]; ok {
				s.EqualValues(dbItem.Quantity, test.patchMap["quantity"])
			}
			if _, ok := test.patchMap["checked"]; ok {
				s.True(dbItem.Checked)
			}
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
