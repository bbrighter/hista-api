package repository

import (
	"time"

	"encore.app/shoppingList/entity"
	"gorm.io/gorm"
)

func (s *RepoTestSuite) TestGetOrCreateMoment() {
	tests := map[string]struct {
		entryExists bool
		expectError error
	}{
		"ok, no entry exists": {},
		"ok, entry exists":    {entryExists: true},
	}

	for name, test := range tests {
		s.Run(name, func() {
			oldTime := time.Date(2000, 10, 5, 3, 2, 1, 0, time.UTC)
			if test.entryExists {
				s.createMoment(oldTime)
			}

			moment, err := s.MomentRepo.GetOrCreate(s.ctx)
			if test.expectError != nil {
				s.ErrorIs(err, test.expectError)
				return
			}
			s.NoError(err)
			if test.entryExists {
				s.True(moment.UpdatedAt.Equal(oldTime), moment.UpdatedAt.String(), oldTime.String())
			} else {
				s.True(time.Until(moment.UpdatedAt).Seconds() < 10)
			}

		})
	}
}

func (s *RepoTestSuite) TestUpdateMoments() {
	oldTime := time.Date(2000, 10, 5, 3, 2, 1, 0, time.UTC)
	tests := map[string]struct {
		useExisting bool
		expectError error
	}{
		"ok":        {useExisting: true},
		"not found": {expectError: gorm.ErrRecordNotFound},
	}

	for name, test := range tests {
		s.Run(name, func() {
			if test.useExisting {
				s.createMoment(oldTime)
			}

			err := s.MomentRepo.Update(s.ctx)
			if test.expectError != nil {
				s.ErrorIs(err, test.expectError)
				return
			}
			s.NoError(err)
			moment, _ := gorm.G[entity.Moment](s.db).First(s.ctx)
			s.True(time.Until(moment.UpdatedAt).Seconds() < 5)
		})
	}
}


func (s *RepoTestSuite) TestCaching() {	
	// No cache exists
	updatedAt := time.Now()
	s.createMoment(updatedAt)

	moment, err := s.MomentRepo.GetOrCreate(s.ctx)
	s.NoError(err)
	s.Equal(updatedAt.Round(time.Second), moment.UpdatedAt.Round(time.Second))

	// Remove db entry, rely on cached value
	rows, err := gorm.G[entity.Moment](s.db).Where("1=1").Delete(s.ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(1, rows)
	
	moment, err = s.MomentRepo.GetOrCreate(s.ctx)
	s.NoError(err)
	s.Equal(updatedAt.Round(time.Second), moment.UpdatedAt.Round(time.Second))


	
}