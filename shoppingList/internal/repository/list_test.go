package repository

import (
	"testing"

	"encore.app/shoppingList/entity"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(ShoppingListTestSuite))
}

func (s *ShoppingListTestSuite) TestFirstOrCreateList() {
	tests := map[string]struct {
		expectedError error
	}{
		"ok": {},
	}

	for name, test := range tests {
		s.Run(name, func() {
			list, err := s.ListRepo.FirstOrCreate(s.ctx)
			if test.expectedError != nil {
				s.ErrorIs(err, test.expectedError)
				return
			}
			s.NoError(err)
			s.EqualValues(1, list.ID)
			sameList, err := s.ListRepo.FirstOrCreate(s.ctx)
			rows, _ := gorm.G[entity.List](s.ListRepo.db).Count(s.ctx, "*")
			s.EqualValues(rows, 1)
			s.NoError(err)
			s.Equal(list, sameList)

			rows, _ = gorm.G[entity.List](s.ListRepo.db).Count(s.ctx, "*")
			s.EqualValues(rows, 1)
		})
	}
}
