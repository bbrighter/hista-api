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

	var cat = &entity.SymptomCategory{Name: "cat"}
	err := repo.CreateCategory(cat)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cat.ID)
}

func TestRenameCategory(t *testing.T) {
	repo := initTest(t)

	var cat = &entity.SymptomCategory{ID: 1, Name: "cat"}
	err := repo.db.Create(cat).Error
	assert.NoError(t, err)

	err = repo.RenameCategory(cat, "new name")
	assert.NoError(t, err)
	assert.Equal(t, "new name", cat.Name)

	var nonExistingCat = &entity.SymptomCategory{ID: 100, Name: "cat"}
	err = repo.RenameCategory(nonExistingCat, "new name")
	assert.Error(t, err)
}

func TestDeleteCategory(t *testing.T) {
	tests := map[string]struct {
		catId         uint
		symptoms      []entity.Symptom
		expectedError string
	}{
		"ok":                      {catId: 1},
		"not found":               {catId: 100, expectedError: "not_found: not found"},
		"symptoms block deleting": {catId: 1, symptoms: []entity.Symptom{{Name: "symptom"}}, expectedError: "failed_precondition: cannot delete category with symptoms"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo := initTest(t)
			err := repo.db.Create(&entity.SymptomCategory{
				ID:       1,
				Symptoms: test.symptoms,
			}).Error
			assert.NoError(t, err)

			err = repo.DeleteCategory(test.catId)
			if test.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})

	}
}
