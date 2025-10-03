package api

import (
	"testing"
	"time"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
)

func TestHeadachesAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	var headaches entity.HeadachesResponse
	headaches, err = service.ListHeadaches(ctx, TEST_PIID)
	assert.NoError(t, err)
	assert.Len(t, headaches.Headaches, 0)

	// Create one headache
	date := time.Date(2018, 1, 2, 3, 4, 5, 0, time.Local)
	var severity entity.HeadacheSeverity = 5
	idResp, err := service.PostHeadache(ctx, TEST_PIID, PostHeadacheParams{Date: date, Severity: severity})
	id := idResp.ID
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)

	// Get headache
	headaches, err = service.ListHeadaches(ctx, TEST_PIID)
	assert.NoError(t, err)
	assert.Len(t, headaches.Headaches, 1)

	headache, err := service.GetHeadache(ctx, TEST_PIID, id)
	assert.NoError(t, err)
	assert.Equal(t, date, headache.Date)
	assert.Equal(t, severity, headache.Severity)
	assert.Nil(t, headache.Positions)

	// Patch and verify
	newDate := time.Date(2019, 1, 2, 3, 4, 5, 0, time.Local)
	err = service.PatchHeadacheDate(ctx, TEST_PIID, id, PatchHeadacheDateParams{Date: newDate})
	assert.NoError(t, err)
	var newSeverity entity.HeadacheSeverity = 1
	err = service.PatchHeadacheSeverity(ctx, TEST_PIID, id, PatchHeadacheSeverityParams{Severity: newSeverity})
	assert.NoError(t, err)
	newTypes := entity.HeadacheTypes{entity.Dull}
	err = service.PatchHeadacheTypes(ctx, TEST_PIID, id, PatchHeadacheTypesParams{Types: newTypes})
	assert.NoError(t, err)
	newPositions := entity.HeadachePositions{entity.Back, entity.Ear}
	err = service.PatchHeadachePositions(ctx, TEST_PIID, id, PatchHeadachePositionsParams{Positions: newPositions})
	assert.NoError(t, err)
	newSymptoms := entity.HeadacheSymptoms{entity.Dizziness}
	err = service.PatchHeadacheSymptoms(ctx, TEST_PIID, id, PatchHeadacheSymptomsParams{Symptoms: newSymptoms})
	assert.NoError(t, err)

	headache, err = service.GetHeadache(ctx, TEST_PIID, id)
	assert.NoError(t, err)
	assert.Equal(t, newDate, headache.Date)
	assert.Equal(t, newPositions, headache.Positions)
	assert.Equal(t, newSeverity, headache.Severity)
	assert.Equal(t, newSymptoms, headache.Symptoms)
	assert.Equal(t, newTypes, headache.Types)

	// Delete and verify
	err = service.DeleteHeadache(ctx, TEST_PIID, idResp.ID)
	assert.NoError(t, err)
	headache, err = service.GetHeadache(ctx, TEST_PIID, idResp.ID)
	assert.Error(t, err)
	headaches, err = service.ListHeadaches(ctx, TEST_PIID)
	assert.NoError(t, err)
	assert.Len(t, headaches.Headaches, 0)
}
