package symptoms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrReplace(t *testing.T) {
	service := initTest(t)

	var err error

	// Verify idempotency
	var existingSymptom = &Symptom{Name: testSymptom.Name, SymptomCategoryID: testCategory.ID}
	err = existingSymptom.createOrReplace(service)
	assert.NoError(t, err)
	assert.Equal(t, existingSymptom.ID, testSymptom.ID)

	// Verify idempotency up to trimming
	var existingSymptomWithWhitespace = &Symptom{Name: " " + testSymptom.Name + " ", SymptomCategoryID: testCategory.ID}
	err = existingSymptomWithWhitespace.createOrReplace(service)
	assert.NoError(t, err)
	assert.Equal(t, existingSymptomWithWhitespace.ID, testSymptom.ID)

	var differentConditonType = &Symptom{Name: "Other symptom", SymptomCategoryID: testCategory.ID}
	err = differentConditonType.createOrReplace(service)
	assert.NoError(t, err)
	assert.NotEqual(t, existingSymptom.ID, differentConditonType.ID)

	// Cleanup
	err = service.db.Delete(&differentConditonType).Error
	assert.NoError(t, err)
}

func TestDeleteIfUnused(t *testing.T) {
	service := initTest(t)

	var unusedSymptom = Symptom{Name: "Unused symptom", SymptomCategoryID: testCategory.ID}
	var err error
	err = service.db.Create(&unusedSymptom).Error
	assert.NoError(t, err)

	err = unusedSymptom.deleteIfUnused(service.db)
	assert.NoError(t, err)

	var symptom = Symptom{ID: unusedSymptom.ID}
	var rows int64 = service.db.First(symptom).RowsAffected
	assert.EqualValues(t, 0, rows)

	var usedSymptom = Symptom{ID: testSymptom.ID}

	err = usedSymptom.deleteIfUnused(service.db)
	assert.NoError(t, err)

	rows = service.db.First(&usedSymptom).RowsAffected
	assert.EqualValues(t, rows, 1)

}
