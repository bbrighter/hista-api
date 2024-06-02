package statistics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDiaryFrom(t *testing.T) {
	t.Parallel()

	var input Input = testInput()
	var diaries []RawDiary = createRawDiary(input)
	assert.Len(t, diaries, 4)
	var firstDiary RawDiary = diaries[0] // First is latest
	assert.Equal(t, firstDiary.Content, "Ingredient")
	assert.Equal(t, firstDiary.Date.Day(), time.Now().Add(2*time.Hour).Day())
	assert.Equal(t, firstDiary.Severity, "cooked")
	var secondDiary RawDiary = diaries[1] // Second happened earlier
	assert.Equal(t, secondDiary.Content, "Symptom")
	assert.Equal(t, secondDiary.Severity, "4")
	assert.Equal(t, secondDiary.Category, "Category")
	var thirdDiary RawDiary = diaries[2]
	assert.Equal(t, "Note", thirdDiary.Content)
	var fourthDiary RawDiary = diaries[3]
	assert.Equal(t, Pollen, fourthDiary.Type)
	assert.Equal(t, "Ambrosia", fourthDiary.Category)
	assert.Equal(t, "Keine bis geringe", fourthDiary.Severity)
}
