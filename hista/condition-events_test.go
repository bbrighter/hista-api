package hista

import (
	"time"

	"encore.app/hista/internal/symptoms"
	"gorm.io/gorm"

	"encore.dev/beta/errs"
	"encore.dev/types/option"
)

func (s *ApiTestSuite) TestCreateConditionEvent() {
	resp, err := s.service.CreateConditionEvent(s.ctx, s.piid)
	s.NoError(err)
	s.Greater(resp.ID, uint(0))
}

func (s *ApiTestSuite) TestListConditionEvents() {
	resp, err := s.service.ListConditionEvents(s.ctx, s.piid)
	s.NoError(err)
	s.Len(resp.ConditionEvents, 0)

	s.createTestEvent()

	resp, err = s.service.ListConditionEvents(s.ctx, s.piid)
	s.NoError(err)
	s.Len(resp.ConditionEvents, 1)
}

func (s *ApiTestSuite) TestGetConditionEvent() {
	_, err := s.service.GetConditionEvent(s.ctx, s.piid, 100)
	s.assertErrCode(err, errs.NotFound)

	id := s.createTestEvent()

	resp, err := s.service.GetConditionEvent(s.ctx, s.piid, id)
	s.NoError(err)
	s.EqualValues(id, resp.ID)
	s.True(time.Now().After(resp.Date))
}

func (s *ApiTestSuite) TestPatchConditionEvent() {
	tests := map[string]struct {
		createTestEvent bool
		date            time.Time
		expectedErrCode errs.ErrCode
	}{
		"not found": {createTestEvent: false, date: time.Now(), expectedErrCode: errs.NotFound},
		"no date":   {createTestEvent: true, expectedErrCode: errs.InvalidArgument},
		"ok":        {createTestEvent: true, date: time.Now()},
	}

	for name, test := range tests {
		s.Run(name, func() {
			var params = ConditionEventRequestParams{Date: test.date}
			var eventId uint = 1000
			if test.createTestEvent {
				id := s.createTestEvent()
				eventId = id
			}
			err := s.service.PatchDate(s.ctx, s.piid, eventId, params)
			s.assertErrCode(err, test.expectedErrCode)
		})

	}
}

func (s *ApiTestSuite) TestDeleteConditionEvent() {
	_, err := s.service.DeleteConditionEvent(s.ctx, s.piid, 10)
	s.assertErrCode(err, errs.NotFound)

	id := s.createTestEvent()
	cats, err := s.service.DeleteConditionEvent(s.ctx, s.piid, id)
	s.NoError(err)
	s.Len(cats.Categories, 0)
}

func (s *ApiTestSuite) TestPostCondition() {
	tests := map[string]struct {
		eventId               uint
		symptomName           bool
		symptomId             bool
		expectErrCode         errs.ErrCode
		expectedSymptomLength int
	}{
		"ok, symptom name": {symptomName: true},
		"ok, use Ids":      {symptomId: true},
		"not found":        {expectErrCode: errs.InvalidArgument, eventId: 1000, symptomName: true},
	}
	for name, test := range tests {
		s.Run(name, func() {
			id := s.createTestEvent()

			eventId := test.eventId
			if test.eventId == 0 {
				eventId = id
			}

			var params ConditionRequestParams
			if test.symptomName {
				idResp, err := s.service.PostSymptomCategory(s.ctx, s.piid, PostSymptomCategoryRequest{Name: "name"})
				s.NoError(err)
				var name string = "name"
				params.SymptomName = option.Some(name)
				params.CategoryID = option.Some(idResp.ID)
			}
			if test.symptomId {
				idResp, err := s.service.PostSymptomCategory(s.ctx, s.piid, PostSymptomCategoryRequest{Name: "name"})
				s.NoError(err)

				params.CategoryID = option.Some(idResp.ID)

				var symptom = symptoms.Symptom{Name: "symptom", SymptomCategoryID: idResp.ID}
				gorm.G[symptoms.Symptom](s.service.DB).Create(s.ctx, &symptom)

				params.SymptomID = option.Some(symptom.ID)
			}

			_, err := s.service.PostCondition(s.ctx, s.piid, eventId, params)

			s.assertErrCode(err, test.expectErrCode)
		})
	}
}
