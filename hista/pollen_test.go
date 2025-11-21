package hista

import (
	"fmt"
)

func (s *ApiTestSuite) TestGetPollens() {
	tests := map[string]struct {
		createPollen  bool
		expectedCount int
	}{
		"0": {},
		"1": {createPollen: true, expectedCount: 1},
	}
	for name, test := range tests {
		s.Run(name, func() {
			if test.createPollen {
				fmt.Print("create pollen")
				s.createTestPollen()
			}
			pollens, err := s.service.ListPollens(s.ctx, s.piid)
			s.NoError(err)
			s.Len(pollens.Pollens, test.expectedCount)
		})
	}
}

func (s *ApiTestSuite) TestUpdatePollen() {
	err := s.service.UpdatePollen(s.ctx)
	s.NoError(err)
}
