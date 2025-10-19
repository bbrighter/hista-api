package repository

import (
	"context"

	"encore.app/shared/contextKeys"
	"encore.app/shared/generic_queries"
	"encore.app/shoppingList/entity"
	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

func (s *ShoppingListTestSuite) TestCreateItem() {
	tests := map[string]struct {
		useNonExistingId bool
		useWrongPiid     bool
		expectError      string
	}{
		"ok":            {},
		"wrong prod id": {useNonExistingId: true, expectError: "23503"},
		"wrong piid":    {useWrongPiid: true, expectError: "23503"},
	}

	for name, test := range tests {
		s.Run(name, func() {
			// _, repo, _, ctx := initTest(t)
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
			itemId, err := s.ItemRepo.Create(createCtx, prodId, list.ID)
			if test.expectError != "" {
				s.AssertPostgresError(err, test.expectError)
				// pgErr, ok := err.(*pgconn.PgError)
				// s.True(ok, "expected Postgres error")
				// s.Equal(test.expectError, pgErr.Code)
			} else {
				s.EqualValues(itemId, 1)
			}
		})
	}
}

// func TestDeleteItem(t *testing.T) {

// }

func (s *ShoppingListTestSuite) TestCheckItem() {
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
			_, item := s.createItem()

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
			respItem, err := generic_queries.First[*entity.Item](ctx, s.ItemRepo.db, itemId)
			s.Require().NoError(err)
			s.Equal(test.expectedValue, respItem.Checked)
		})
	}
}

func (s *ShoppingListTestSuite) TestListItems() {
	tests := map[string]struct {
		useWrongPiid bool
		expectedLen  int
	}{
		"ok":         {expectedLen: 1},
		"wrong piid": {useWrongPiid: true, expectedLen: 0},
	}

	for name, test := range tests {
		s.Run(name, func() {
			s.createItem()

			ctx := s.ctx
			if test.useWrongPiid {
				guid, _ := uuid.NewV4()
				ctx = context.WithValue(ctx, contextKeys.Piid, guid)
			}

			items, err := s.ItemRepo.List(ctx)
			s.NoError(err)
			s.Len(items, test.expectedLen)
		})
	}
}
