package repository

import (
	"context"

	"encore.app/shared/contextKeys"
	"encore.dev/types/uuid"
)

func (s *RepoTestSuite) TestCreate() {
	tests := map[string]struct {
		useExistingName       bool
		createInDifferentPiid bool
		expectErrorCode       string
	}{
		"ok":                                   {},
		"name already exists":                  {useExistingName: true, expectErrorCode: "23505"},
		"name exists, but in a different piid": {createInDifferentPiid: true, useExistingName: true},
	}

	for name, test := range tests {
		s.Run(name, func() {
			initCtx := s.ctx
			if test.createInDifferentPiid {
				newGuid, _ := uuid.NewV4()
				initCtx = context.WithValue(initCtx, contextKeys.Piid, newGuid)
			}
			s.createProduct(initCtx)

			name = "new name"
			if test.useExistingName {
				name = "name"
			}
			_, err := s.ProductRepo.Create(s.ctx, name)
			if test.expectErrorCode != "" {
				s.AssertPostgresError(err, test.expectErrorCode)
				// s.ErrorContains(err, test.expectErrorCode)
			} else {
				s.NoError(err)
			}
		})
	}
}
