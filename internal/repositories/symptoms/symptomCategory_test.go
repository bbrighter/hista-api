package symptoms

import (
	"testing"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
)

func TestListCategories(t *testing.T) {
	repo := initTest(t)

	cats := repo.ListCategories()
	assert.Len(t, cats, 0)

	repo.db.Create(&entity.SymptomCategory{
		Name: "cat",
		Symptoms: []entity.Symptom{
			{Name: "symptom"},
		},
	})
	cats = repo.ListCategories()
	assert.Len(t, cats, 1)
	assert.Equal(t, "cat", cats[0].Name)
	assert.Len(t, cats[0].Symptoms, 1)
	assert.Equal(t, "symptom", cats[0].Symptoms[0].Name)
}

func TestCreateCategory(t *testing.T) {
	repo := initTest(t)

	id, err := repo.CreateCategory("cat")
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)
}
