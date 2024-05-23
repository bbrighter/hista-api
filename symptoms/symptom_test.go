package symptoms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrReplace(t *testing.T) {
	service := initTest(t)

	var err error

	var symptom = &Symptom{Name: "Symptom", SymptomCategoryID: 1}
	err = symptom.createOrReplace(service)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, symptom.ID, uint(1))

	var differentConditonType = &Symptom{Name: "Other symptom", SymptomCategoryID: 1}
	err = differentConditonType.createOrReplace(service)
	assert.NoError(t, err)
	assert.NotEqual(t, symptom.ID, differentConditonType.ID)

	var sameConditionType = &Symptom{Name: "Symptom", SymptomCategoryID: 1}
	err = sameConditionType.createOrReplace(service)
	assert.NoError(t, err)
	assert.Equal(t, symptom.ID, sameConditionType.ID)

	// Cleanup
	err = service.db.Delete(&differentConditonType).Error
	assert.NoError(t, err)
}

func TestDeleteIfUnused(t *testing.T) {
	service := initTest(t)

	var unusedSymptom = Symptom{ID: 2, Name: "Unused symptom", SymptomCategoryID: 1}
	var err error
	err = service.db.Create(&unusedSymptom).Error
	assert.NoError(t, err)

	err = unusedSymptom.deleteIfUnused(service.db)
	assert.NoError(t, err)

	var symptoms []Symptom
	var rows int64 = service.db.Find(&symptoms).RowsAffected
	assert.EqualValues(t, 1, rows)

	var usedSymptom = Symptom{ID: 1}

	err = usedSymptom.deleteIfUnused(service.db)
	assert.NoError(t, err)

	rows = service.db.Find(&symptoms).RowsAffected
	assert.EqualValues(t, 1, rows)

}
