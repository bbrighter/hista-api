package hista

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"

	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func (service *Service) createTestEvent(ctx context.Context, t *testing.T) uint {
	resp, err := service.CreateConditionEvent(ctx, TEST_PIID)
	require.NoError(t, err)
	return resp.ID
}

func TestCreateConditionEvent(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, err := service.CreateConditionEvent(ctx, TEST_PIID)
	require.NoError(t, err)
	defer service.DeleteConditionEvent(ctx, TEST_PIID, resp.ID)
	assert.EqualValues(t, 1, resp.ID)
}

func TestGetConditionEvents(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, err := service.ListConditionEvents(ctx, TEST_PIID)
	assert.NoError(t, err)
	assert.Len(t, resp.ConditionEvents, 0)

	service.createTestEvent(ctx, t)

	resp, err = service.ListConditionEvents(ctx, TEST_PIID)
	assert.NoError(t, err)
	assert.Len(t, resp.ConditionEvents, 1)
}

func TestGetConditionEvent(t *testing.T) {
	service, ctx := initAPITest(t)

	_, err := service.GetConditionEvent(ctx, TEST_PIID, 100)
	assert.EqualError(t, err, "not_found: not found")

	id := service.createTestEvent(ctx, t)

	resp, err := service.GetConditionEvent(ctx, TEST_PIID, id)
	assert.NoError(t, err)
	assert.EqualValues(t, id, resp.ID)
	assert.True(t, time.Now().After(resp.Date))
}

func TestPatchConditionEvent(t *testing.T) {
	tests := map[string]struct {
		createTestEvent   bool
		date              time.Time
		expectedErrorText string
	}{
		"not found": {createTestEvent: false, date: time.Now(), expectedErrorText: "not_found: not found"},
		"no date":   {createTestEvent: true, expectedErrorText: "invalid_argument: date must be set"},
		"ok":        {createTestEvent: true, date: time.Now()},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			service, ctx := initAPITest(t)
			var params = ConditionEventRequestParams{Date: test.date}
			var eventId uint = 1000
			if test.createTestEvent {
				id := service.createTestEvent(ctx, t)
				eventId = id
			}
			err := service.PatchDate(ctx, TEST_PIID, eventId, params)
			if test.expectedErrorText != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expectedErrorText)
			} else {
				assert.NoError(t, err)
			}
		})

	}
}

func TestDeleteConditionEvent(t *testing.T) {
	service, ctx := initAPITest(t)

	_, err := service.DeleteConditionEvent(ctx, TEST_PIID, 10)
	assert.EqualError(t, err, "not_found: not found")

	id := service.createTestEvent(ctx, t)
	cats, err := service.DeleteConditionEvent(ctx, TEST_PIID, id)
	assert.NoError(t, err)
	assert.Len(t, cats.Categories, 0)
}

func TestPostCondition(t *testing.T) {
	tests := map[string]struct {
		eventId               uint
		symptomName           bool
		symptomId             bool
		expectErrorMsg        string
		expectedSymptomLength int
	}{
		"ok, symptom name": {symptomName: true},
		"ok, use Ids":      {symptomId: true},
		"not found":        {expectErrorMsg: "not_found: not found", eventId: 1000, symptomName: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			service, ctx := initAPITest(t)
			id := service.createTestEvent(ctx, t)

			eventId := test.eventId
			if test.eventId == 0 {
				eventId = id
			}

			var params ConditionRequestParams
			if test.symptomName {
				idResp, err := service.PostSymptomCategory(ctx, TEST_PIID, PostSymptomCategoryRequest{Name: "name"})
				assert.NoError(t, err)
				defer service.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.Symptom{})
				defer service.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.SymptomCategories{})
				var name string = "name"
				params.SymptomName = &name
				params.CategoryID = &idResp.ID
			}
			if test.symptomId {
				idResp, err := service.PostSymptomCategory(ctx, TEST_PIID, PostSymptomCategoryRequest{Name: "name"})
				assert.NoError(t, err)

				params.CategoryID = &idResp.ID

				var symptom = entity.Symptom{Name: "symptom", SymptomCategoryID: idResp.ID, PIID: uuid.FromStringOrNil(TEST_PIID_STR)}
				generic_queries.Create(ctx, service.DB, &symptom)

				defer service.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.Symptom{})
				defer service.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.SymptomCategories{})
				params.SymptomID = &symptom.ID
			}

			_, err := service.PostCondition(ctx, TEST_PIID, eventId, params)

			if test.expectErrorMsg != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expectErrorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
