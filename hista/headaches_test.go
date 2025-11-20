package hista

import "encore.dev/beta/errs"

func (s *ApiTestSuite) TestListHeadaches() {
	tests := map[string]struct {
		createHeadache bool
		expectedAmount int
	}{
		"1": {createHeadache: true, expectedAmount: 1},
		"0": {},
	}
	for name, test := range tests {
		s.Run(name, func() {
			if test.createHeadache {
				s.createTestHeadache()
			}
			resp, err := s.service.ListHeadaches(s.ctx, s.piid)
			s.NoError(err)
			s.Len(resp.Headaches, test.expectedAmount)
		})
	}
}

func (s *ApiTestSuite) TestGetHeadache() {
	tests := map[string]struct {
		useWrongId      bool
		expectedErrCode errs.ErrCode
	}{
		"ok":        {},
		"not found": {useWrongId: true, expectedErrCode: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			id := s.createTestHeadache()
			if test.useWrongId {
				id = 1000
			}
			resp, err := s.service.GetHeadache(s.ctx, s.piid, id)
			if s.assertErrCode(err, test.expectedErrCode) {
				return
			}
			s.Equal(id, resp.ID)
		})
	}
}

func (s *ApiTestSuite) TestCreateHeadache() {
	s.T().Skip() // Will be tested in integration test
}

func (s *ApiTestSuite) TestDeleteHeadache() {
	tests := map[string]struct {
		useWrongId      bool
		expectedErrCode errs.ErrCode
	}{
		"ok":        {},
		"not found": {useWrongId: true, expectedErrCode: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			id := s.createTestHeadache()
			if test.useWrongId {
				id = 1000
			}
			err := s.service.DeleteHeadache(s.ctx, s.piid, id)
			s.assertErrCode(err, test.expectedErrCode)
		})
	}
}

func (s *ApiTestSuite) TestPatchHeadaches() {
	s.T().Skip() // Will be tested in integration tests
}
