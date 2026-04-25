package hista

import (
	"time"

	"encore.dev/types/option"
)

func (s *ApiTestSuite) TestCreateStatus() {
	resp, err := s.service.PostStatus(s.ctx, s.piid, DateParam{time.Now()})
	s.NoError(err)
	s.GreaterOrEqual(resp.ID, uint(1))
}

// func TestAddStatus(t *testing.T) {
// 	service, ctx := initAPITest(t)

// 	var err error
// 	var resp, morningResp, eveningResp entity.StatusResponse
// 	resp, err = service.PostStatus(ctx, DateParam{time.Now()})
// 	defer service.DeleteStatus(ctx, resp.ID)
// 	assert.NoError(t, err)

// 	params := StatusParams{
// 		TimeOfDay: entity.Morning,
// 		Fitness:   entity.Good,
// 		Sleep:     entity.Good,
// 	}
// 	morningResp, err = service.PutStatus(ctx, resp.ID, params)
// 	assert.NoError(t, err)
// 	assert.GreaterOrEqual(t, morningResp.Morning.ID, uint(1))

// 	// Error if adding second time without fitting MorningID
// 	_, err = service.PutStatus(ctx, resp.ID, params)
// 	assert.Error(t, err)

// 	eveningParams := StatusParams{
// 		TimeOfDay: entity.Evening,
// 		Fitness:   entity.Middle,
// 	}
// 	eveningResp, err = service.PutStatus(ctx, resp.ID, eveningParams)
// 	assert.NoError(t, err)
// 	assert.Equal(t, entity.Middle, eveningResp.Evening.Fitness)
// 	assert.Equal(t, entity.Good, eveningResp.Morning.Fitness)

// }

// func TestAddStatusIdempotent(t *testing.T) {
// 	service, ctx := initAPITest(t)

// 	resp, err := service.PostStatus(ctx, DateParam{time.Now()})
// 	defer service.DeleteStatus(ctx, resp.ID)
// 	assert.NoError(t, err)

// 	params := StatusParams{
// 		TimeOfDay: entity.Morning,
// 		Fitness:   entity.Good,
// 		Sleep:     entity.Good,
// 	}
// 	statusResp, err := service.PutStatus(ctx, resp.ID, params)
// 	assert.NoError(t, err)
// 	assert.GreaterOrEqual(t, statusResp.Morning.ID, uint(1))
// 	params.ID = statusResp.Morning.ID
// 	newStatusResp, err := service.PutStatus(ctx, resp.ID, params)
// 	assert.NoError(t, err)
// 	assert.Equal(t, newStatusResp.Morning.ID, statusResp.Morning.ID)
// }

func (s *ApiTestSuite) TestListStatus() {
	tests := map[string]struct {
		createStatus   bool
		expectedAmount int
	}{
		"0": {},
		"1": {createStatus: true, expectedAmount: 1},
	}
	for name, test := range tests {
		s.Run(name, func() {
			if test.createStatus {
				s.createTestStatus()
			}
			resp, err := s.service.ListStatus(s.ctx, s.piid)
			s.NoError(err)
			s.Len(resp.Statuses, test.expectedAmount)
		})
	}
}

func (s *ApiTestSuite) TestUpdateStatus() {
	tests := map[string]struct {
		date    bool
		morning bool
		evening bool
	}{
		"only date":    {date: true},
		"only morning": {morning: true},
		"only evening": {evening: true},
	}

	for name, test := range tests {
		s.Run(name, func() {
			statusId := s.createTestStatus()

			var params = PatchStatusParams{}
			if test.date {
				params.Date = option.Some(time.Date(2022, 6, 5, 4, 3, 2, 0, time.UTC))
			}
			morningFitness := 1
			if test.morning {
				params.MorningFitness = option.Some(morningFitness)
			}
			eveningFitness := 3
			if test.evening {
				params.EveningFitness = option.Some(eveningFitness)
			}

			err := s.service.PatchStatus(s.ctx, s.piid, statusId, params)
			s.NoError(err)

			statuses, err := s.service.status.ListStatuses(s.ctx)
			s.NoError(err)
			status := statuses[0]
			if test.date {
				s.True(time.Date(2022, 6, 5, 4, 3, 2, 0, time.UTC).Equal(status.Date))
			} else {
				s.True(time.Now().After(status.Date))
			}
			if test.morning {
				s.EqualValues(&morningFitness, status.MorningFitness)
			}
			if test.evening {
				s.EqualValues(&eveningFitness, status.EveningFitness)
			}
		})
	}

}
