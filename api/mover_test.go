package api

import (
	"context"
	"testing"

	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMove(t *testing.T) {
	newPiid := uuid.FromStringOrNil("09d5289f-d158-4dd8-b369-c4d58110e921")
	s, ctx := initAPITest(t)

	s.createTestEvent(ctx, t)
	s.createTestMeal(ctx, t)
	s.createTestNote(ctx, t)

	err := s.MovePiid(ctx, uuid.FromStringOrNil(TEST_PIID_STR), newPiid)
	assert.NoError(t, err)

	condResp, _ := s.GetConditionEvents(ctx)
	assert.Len(t, condResp.ConditionEvents, 0)
	mealResp, _ := s.GetMeals(ctx)
	assert.Len(t, mealResp.Meals, 0)
	notesResp, _ := s.GetNotes(ctx)
	assert.Len(t, notesResp.Notes, 0)

	newCtx := context.WithValue(ctx, "piid", newPiid)
	condResp, _ = s.GetConditionEvents(newCtx)
	assert.Len(t, condResp.ConditionEvents, 1)
	mealResp, _ = s.GetMeals(newCtx)
	assert.Len(t, mealResp.Meals, 1)
	notesResp, _ = s.GetNotes(newCtx)
	assert.Len(t, notesResp.Notes, 1)
}
