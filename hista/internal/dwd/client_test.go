package dwd

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestData() (dwd DWD, err error) {
	var content []byte
	content, err = os.ReadFile("./mockResponse.json")
	if err != nil {
		return DWD{}, err
	}
	err = json.Unmarshal(content, &dwd)
	return dwd, err
}

func TestGetKarlsruheDate(t *testing.T) {
	client := &DWDClient{query: getTestData}

	dwd, updatedAt, err := client.GetKarlsruheData()
	assert.NoError(t, err)

	berlin, err := time.LoadLocation("Europe/Berlin")
	require.NoError(t, err)
	expectedTime := time.Date(2024, 5, 31, 11, 0, 0, 0, berlin)
	assert.Equal(t, expectedTime, updatedAt)
	assert.Equal(t, "0", dwd.Hasel.Today)
	assert.Equal(t, "0", dwd.Hasel.Today)
	assert.Equal(t, "0", dwd.Esche.Today)
	assert.Equal(t, "1", dwd.Graeser.Today)
	assert.Equal(t, "2", dwd.Ambrosia.Today)
	assert.Equal(t, "1-2", dwd.Erle.Today)
	assert.Equal(t, "0-1", dwd.Roggen.Today)
	assert.Equal(t, "3", dwd.Birke.Today)
	assert.Equal(t, "2-3", dwd.Beifuss.Today)
}
