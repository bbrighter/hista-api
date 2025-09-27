package symptoms

import (
	"testing"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
)

func TestListCategories(t *testing.T) {
	repo, ctx := initTest(t)

	cats, err := repo.ListCategories(ctx)
	assert.NoError(t, err)
	assert.Len(t, cats, 0)

	repo.db.Create(&entity.SymptomCategory{
		Name: "cat",
		PIID: GUID,
		Symptoms: []entity.Symptom{
			{Name: "symptom", PIID: GUID},
		},
	})
	cats, err = repo.ListCategories(ctx)
	assert.NoError(t, err)
	assert.Len(t, cats, 1)
	assert.Equal(t, "cat", cats[0].Name)
	assert.Len(t, cats[0].Symptoms, 1)
	assert.Equal(t, "symptom", cats[0].Symptoms[0].Name)
}

func TestCreateCategory(t *testing.T) {
	tests := map[string]struct {
		catNameSet           bool
		catNameExistsAlready bool
		expectError          bool
	}{
		"ok":           {catNameSet: true},
		"ok and exits": {catNameSet: true, catNameExistsAlready: true},
		"name missing": {catNameSet: false, expectError: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, ctx := initTest(t)
			if test.catNameExistsAlready {
				repo.CreateCategory(ctx, &entity.SymptomCategory{Name: "cat"})
			}
			var cat = &entity.SymptomCategory{}
			if test.catNameSet {
				cat.Name = "cat"
			}

			err := repo.CreateCategory(ctx, cat)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Greater(t, cat.ID, uint(0))
			}
		})
	}
}

func TestRenameCategory(t *testing.T) {
	repo, ctx := initTest(t)

	var cat = &entity.SymptomCategory{ID: 1, Name: "cat", PIID: GUID}
	err := repo.db.Create(cat).Error
	assert.NoError(t, err)

	err = repo.RenameCategory(ctx, cat, "new name")
	assert.NoError(t, err)
	assert.Equal(t, "new name", cat.Name)

	var nonExistingCat = &entity.SymptomCategory{ID: 100, Name: "cat", PIID: GUID}
	err = repo.RenameCategory(ctx, nonExistingCat, "new name")
	assert.Error(t, err)
}

func TestDeleteCategory(t *testing.T) {
	tests := map[string]struct {
		catId         uint
		symptoms      []entity.Symptom
		expectedError string
	}{
		"ok":                      {catId: 1},
		"not found":               {catId: 100, expectedError: "record not found"},
		"symptoms block deleting": {catId: 1, symptoms: []entity.Symptom{{Name: "symptom", PIID: GUID}}, expectedError: "cannot delete category with symptoms"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, ctx := initTest(t)
			err := repo.db.Create(&entity.SymptomCategory{
				ID:       1,
				Symptoms: test.symptoms,
				PIID:     GUID,
			}).Error
			assert.NoError(t, err)

			err = repo.DeleteCategory(ctx, test.catId)
			if test.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})

	}
}
