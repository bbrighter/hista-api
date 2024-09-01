package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDiary(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, err := service.GetDiary(ctx)
	assert.NoError(t, err)
	assert.Len(t, resp.Diaries, 0)

	symptomCleanup := service.createTestSymptom(t)
	defer symptomCleanup(t)

	cleanupFood := service.createTestFood(t)
	defer cleanupFood(t)

	cleanupNote := service.createTestNote(t)
	defer cleanupNote(t)

	resp, err = service.GetDiary(ctx)
	assert.NoError(t, err)
	assert.Len(t, resp.Diaries, 3)
}
