package symptoms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSymptomsAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var resp SymptomCategoriesResponse
	var err error
	resp, err = service.GetSymptoms(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(resp.Categories), 1)
}
