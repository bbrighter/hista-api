package hista

import (
	"encore.app/hista/entity"
	"encore.dev/beta/errs"
)

func (s *ApiTestSuite) TestPatchCondition() {
	s.T().Skip()

	id := s.createTestEvent()

	var err error
	err = s.service.PatchCondition(s.ctx, s.piid, 100, PatchSeverityRequestParams{Severity: entity.HighSeverity})
	s.assertErrCode(err, errs.NotFound)

	catId, err := s.service.symptoms.CreateCategory(s.ctx, "cat")
	s.NoError(err)
	symtpomName := "name"
	resp, err := s.service.PostCondition(s.ctx, s.piid, id, ConditionRequestParams{SymptomName: &symtpomName, CategoryID: &catId})
	s.NoError(err)
	conditionId := resp.Condition.ID

	err = s.service.PatchCondition(s.ctx, s.piid, conditionId, PatchSeverityRequestParams{Severity: entity.HighSeverity})
	s.NoError(err)
}

func (s *ApiTestSuite) TestDeleteCondition() {
	tests := map[string]struct {
		useWrongId      bool
		expectedErrCode errs.ErrCode
	}{
		"ok":        {},
		"not found": {useWrongId: true, expectedErrCode: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			condId, _, _ := s.createTestCondition()
			if test.useWrongId {
				condId = 1000
			}
			cats, err := s.service.DeleteCondition(s.ctx, s.piid, condId)
			s.assertErrCode(err, test.expectedErrCode)
			if test.expectedErrCode == 0 {
				s.Len(cats.Categories, 1, "categories are not deleted")
			}

		})
	}
}
