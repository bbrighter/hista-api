package symptoms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSymptomsAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)

	var resp SymptomCategoriesResponse
	var err error
	resp, err = service.GetSymptoms(ctx)
	assert.NoError(t, err)
	assert.Len(t, resp.Categories, 0)

	service.testCreateConditionEvent(t)
	resp, err = service.GetSymptoms(ctx)
	assert.NoError(t, err)
	assert.Len(t, resp.Categories, 1)
}
