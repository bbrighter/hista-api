package hista

import (
	"context"

	"encore.app/shared/contextKeys"
	"encore.dev/types/uuid"
)

func (s *ApiTestSuite) TestMove() {
	newPiid := uuid.FromStringOrNil("09d5289f-d158-4dd8-b369-c4d58110e921")

	s.createTestEvent()
	s.createTestMeal()
	s.createTestNote()

	err := s.service.MovePiid(s.ctx, s.piid, newPiid)
	s.NoError(err)

	condResp, _ := s.service.ListConditionEvents(s.ctx, s.piid)
	s.Len(condResp.ConditionEvents, 0)
	mealResp, _ := s.service.ListMeals(s.ctx, s.piid)
	s.Len(mealResp.Meals, 0)
	notesResp, _ := s.service.ListNotes(s.ctx, s.piid)
	s.Len(notesResp.Notes, 0)

	newCtx := context.WithValue(s.ctx, contextKeys.Piid, newPiid)
	condResp, _ = s.service.ListConditionEvents(newCtx, newPiid)
	s.Len(condResp.ConditionEvents, 1)
	mealResp, _ = s.service.ListMeals(newCtx, newPiid)
	s.Len(mealResp.Meals, 1)
	notesResp, _ = s.service.ListNotes(newCtx, newPiid)
	s.Len(notesResp.Notes, 1)
}
