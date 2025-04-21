package symptoms

import (
	"testing"

	"encore.app/entity"
	"encore.app/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrReplace(t *testing.T) {
	repo := initTest(t)

	id, err := repo.CreateOrReplace("new name", 1)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)

	id, err = repo.CreateOrReplace("new name 2", 1)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, id)

	id, err = repo.CreateOrReplace("new name", 2)
	assert.NoError(t, err)
	assert.EqualValues(t, 3, id)

	id, err = repo.CreateOrReplace("new name", 1)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)
}

func TestChangeCategory(t *testing.T) {
	tests := map[string]struct {
		symptomId        uint
		targetCategoryId uint
		expectedError    error
	}{
		"ok":          {symptomId: 1, targetCategoryId: 1},
		"no symptom":  {symptomId: 100, targetCategoryId: 1, expectedError: errors.ErrorNotFound},
		"no category": {symptomId: 1, targetCategoryId: 100, expectedError: errors.ErrorNotFound},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo := initTest(t)
			err := repo.db.Create(&entity.SymptomCategory{
				ID:       10,
				Symptoms: []entity.Symptom{{ID: 1}}},
			).Error
			assert.NoError(t, err)
			err = repo.db.Create(&entity.SymptomCategory{ID: 1}).Error
			assert.NoError(t, err)

			err = repo.ChangeCategory(test.symptomId, test.targetCategoryId)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				assert.NoError(t, err)
				var newSymptom entity.Symptom
				repo.db.Take(&newSymptom, test.symptomId)
				assert.EqualValues(t, test.targetCategoryId, newSymptom.SymptomCategoryID)
			}
		})
	}
}

func TestRenameSymptom(t *testing.T) {
	tests := map[string]struct {
		symptomId     uint
		newName       string
		expectedError error
	}{
		"ok":        {symptomId: 1, newName: "new name"},
		"not found": {symptomId: 10, newName: "new name", expectedError: errors.ErrorNotFound},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo := initTest(t)
			err := repo.db.Create(&entity.Symptom{
				ID:   1,
				Name: "old name",
			}).Error
			assert.NoError(t, err)

			err = repo.RenameSymptom(test.symptomId, test.newName)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				assert.NoError(t, err)
				var newSymptom entity.Symptom
				repo.db.Take(&newSymptom, test.symptomId)
				assert.EqualValues(t, test.newName, newSymptom.Name)
			}
		})
	}
}
