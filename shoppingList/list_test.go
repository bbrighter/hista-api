package shoppinglist

import (
	"testing"
	"time"

	sl "encore.app/shoppingList/internal/shoppingList"
	"encore.dev/beta/errs"
	"encore.dev/types/option"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *ApiTestSuite) TestDeleteList() {
	tests := map[string]struct {
		useWrongId        bool
		useWrongPiid      bool
		statusCode        errs.ErrCode
		hasUncheckedItems bool
		forceDelete       bool
		deleteParams      option.Option[string]
	}{
		"ok":               {},
		"not found":        {useWrongId: true, statusCode: errs.NotFound},
		"has items left":   {hasUncheckedItems: true, statusCode: errs.InvalidArgument},
		"wrong piid":       {useWrongPiid: true, statusCode: errs.NotFound},
		"old force delete": {hasUncheckedItems: true, forceDelete: true},
		"force delete":     {hasUncheckedItems: true, deleteParams: option.Some("force")},
		"move delete":      {hasUncheckedItems: true, deleteParams: option.Some("move")},
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
			err := s.service.DeleteList(ctx, s.piid, listId, DeleteListForceDeleteParam{Force: test.forceDelete, DeleteOption: test.deleteParams})
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

func (s *ApiTestSuite) TestMoveDeleteList_checkResults() {
	listId := s.createList()
	s.createItem(listId)

	err := s.service.DeleteList(s.ctx, s.piid, listId, DeleteListForceDeleteParam{DeleteOption: option.Some("move")})
	s.NoError(err)
	list, err := s.service.PostOrGetList(s.ctx, s.piid)
	s.Require().NoError(err)
	s.Len(list.Items, 1)
}

func (s *ApiTestSuite) TestForceDeleteList_checkResults() {
	listId := s.createList()
	s.createItem(listId)

	err := s.service.DeleteList(s.ctx, s.piid, listId, DeleteListForceDeleteParam{DeleteOption: option.Some("force")})
	s.NoError(err)
	list, err := s.service.PostOrGetList(s.ctx, s.piid)
	s.Require().NoError(err)
	s.Len(list.Items, 0)
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

func TestToListResponse(t *testing.T) {
	created := time.Now()
	beforeCreated := time.Now().Add(-time.Minute)

	list := sl.List{ID: 3, PIID: uuid.UUID{}, Items: []sl.Item{
		{ID: 1, ProductId: 3, Checked: false, CreatedAt: created},
		{ID: 2, ProductId: 4, Checked: true, CreatedAt: beforeCreated},
	}}

	resp := toListResponse(list)

	require.Len(t, resp.Items, 2)
	item1 := resp.Items[0]
	assert.EqualValues(t, item1.ID, 1)
	item2 := resp.Items[1]
	assert.EqualValues(t, item2.ID, 2)

}
