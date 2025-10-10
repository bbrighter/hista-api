package symptoms

import (
	"testing"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateOrReplace(t *testing.T) {
	repo, ctx := initTest(t)

	id, err := repo.CreateOrReplace(ctx, "new name", 1)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)

	id, err = repo.CreateOrReplace(ctx, "new name 2", 1)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, id)

	id, err = repo.CreateOrReplace(ctx, "new name", 2)
	assert.NoError(t, err)
	assert.EqualValues(t, 3, id)

	id, err = repo.CreateOrReplace(ctx, "new name", 1)
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
		"no symptom":  {symptomId: 100, targetCategoryId: 1, expectedError: gorm.ErrRecordNotFound},
		"no category": {symptomId: 1, targetCategoryId: 100, expectedError: gorm.ErrRecordNotFound},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, ctx := initTest(t)
			var cat = entity.SymptomCategory{ID: 1, PIID: GUID}
			err := gorm.G[entity.SymptomCategory](repo.db).Create(ctx, &cat)
			require.NoError(t, err)
			var sym = entity.Symptom{ID: 1, PIID: GUID, SymptomCategoryID: cat.ID, SymptomCategoryPIID: GUID}
			err = gorm.G[entity.Symptom](repo.db).Create(ctx, &sym)
			require.NoError(t, err)

			err = repo.ChangeCategory(ctx, test.symptomId, test.targetCategoryId)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.expectedError)
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
		"not found": {symptomId: 10, newName: "new name", expectedError: gorm.ErrRecordNotFound},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, ctx := initTest(t)
			err := repo.db.Create(&entity.Symptom{
				ID:   1,
				Name: "old name",
				PIID: GUID,
			}).Error
			assert.NoError(t, err)

			err = repo.RenameSymptom(ctx, test.symptomId, test.newName)

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
