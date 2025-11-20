package hista

import "encore.dev/beta/errs"

func (s *ApiTestSuite) TestDeleteFood() {
	tests := map[string]struct {
		useWrongId      bool
		expectedErrCode errs.ErrCode
	}{
		"ok":        {},
		"not found": {useWrongId: true, expectedErrCode: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			foodId, _ := s.createTestFood()
			if test.useWrongId {
				foodId = 1000
			}
			ing, err := s.service.DeleteFood(s.ctx, s.piid, foodId)
			s.assertErrCode(err, test.expectedErrCode)
			s.Len(ing.Ingredients, 0)
		})
	}
}

func (s *ApiTestSuite) TestPatchFoodCondition() {
	var params = FoodConditionParams{Condition: "raw"}
	tests := map[string]struct {
		useWrongId      bool
		params          FoodConditionParams
		expectedErrCode errs.ErrCode
	}{
		"ok":                {params: params},
		"not found":         {params: params, useWrongId: true, expectedErrCode: errs.NotFound},
		"no params":         {expectedErrCode: errs.InvalidArgument},
		"invalid condition": {params: FoodConditionParams{Condition: "invalid"}, expectedErrCode: errs.InvalidArgument},
	}
	for name, test := range tests {
		s.Run(name, func() {
			foodId, _ := s.createTestFood()
			if test.useWrongId {
				foodId = 1000
			}
			err := s.service.PatchFoodCondition(s.ctx, s.piid, foodId, test.params)
			s.assertErrCode(err, test.expectedErrCode)
		})
	}
}
