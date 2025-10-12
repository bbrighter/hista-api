package statistics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindSymptomsForFoods(t *testing.T) {
	t.Skip() // Doesn't work for SQLite, has to be tested in integration tests
	repo, ctx := initTest(t)

	resp, err := repo.FindSymptomsForFoods(ctx, time.Now(), time.Now().Add(time.Minute), []uint{1})
	assert.NoError(t, err)
	assert.Len(t, resp, 0)
}

func TestCountSymptoms(t *testing.T) {
	repo, ctx := initTest(t)

	results, err := repo.CountSymptoms(ctx, []uint{1})
	assert.NoError(t, err)
	assert.Len(t, results, 0)

	// TODO: test with data filled
}
