package hista

import (
	"testing"

	"encore.app/hista/entity"
	"encore.dev/beta/errs"
	"github.com/stretchr/testify/assert"
)

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
	var params = FoodConditionParams{Condition: entity.Raw}
	tests := map[string]struct {
		useWrongId      bool
		params          FoodConditionParams
		expectedErrCode errs.ErrCode
	}{
		"ok":        {params: params},
		"not found": {params: params, useWrongId: true, expectedErrCode: errs.NotFound},
		// "no params": {expectedErrCode: errs.InvalidArgument},
		// "invalid condition": {params: FoodConditionParams{Condition: "invalid"}, expectedErrCode: errs.InvalidArgument},
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

func TestFoodParamValidation(t *testing.T) {
	tests := map[string]struct {
		condition     string
		expectedError bool
	}{}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			params := FoodConditionParams{Condition: entity.FoodCondition(test.condition)}
			err := params.Validate()
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})

	}
}
