package states

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrReplace(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var category = SymptomCategory{Name: "Category", ID: 2}
	var err error
	err = service.db.Create(&category).Error
	assert.NoError(t, err)

	var conditonType = &Symptom{Name: "New name", SymptomCategoryID: 2}
	err = conditonType.createOrReplace(service)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, conditonType.ID, uint(1))

	var differentConditonType = &Symptom{Name: "Other name", SymptomCategoryID: 2}
	err = differentConditonType.createOrReplace(service)
	assert.NoError(t, err)
	assert.NotEqual(t, conditonType.ID, differentConditonType.ID)

	var sameConditionType = &Symptom{Name: "New name", SymptomCategoryID: 2}
	err = sameConditionType.createOrReplace(service)
	assert.NoError(t, err)
	assert.Equal(t, conditonType.ID, sameConditionType.ID)
}

func TestDeleteIfUnused(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var category = SymptomCategory{ID: 2, Name: "Other category"}
	var unusedConditionType = Symptom{ID: 1, Name: "Name", SymptomCategoryID: 2}
	var err error
	err = service.db.Create(&category).Error
	assert.NoError(t, err)
	err = service.db.Create(&unusedConditionType).Error
	assert.NoError(t, err)

	err = unusedConditionType.deleteIfUnused(service.db)
	assert.NoError(t, err)

	var conditionTypes []Symptom
	var rows int64 = service.db.Find(&conditionTypes).RowsAffected
	assert.EqualValues(t, 0, rows)

	var state = service.testCreateState(t)
	var usedConditionType Symptom = state.Conditions[0].Symptom

	err = usedConditionType.deleteIfUnused(service.db)
	assert.NoError(t, err)

	rows = service.db.Find(&conditionTypes).RowsAffected
	assert.EqualValues(t, 1, rows)

}
