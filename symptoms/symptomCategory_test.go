package symptoms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSymptomCategories(t *testing.T) {
	service := initTest(t)

	var cats []SymptomCategory = getSymptomCategories(service)
	assert.GreaterOrEqual(t, len(cats), 1)
	var cat1 SymptomCategory = cats[0]
	assert.GreaterOrEqual(t, len(cat1.Symptoms), 1)
	assert.NotEqualValues(t, 0, cat1.Symptoms[0].SymptomCategoryID, "!= 0 means value exists")
}

func TestCreateSymptomCategory(t *testing.T) {
	service := initTest(t)

	var cat = &SymptomCategory{Name: "New Category"}
	var err error
	err = cat.create(service)

	assert.NoError(t, err)
	assert.NotEqualValues(t, cat.ID, 0, "ID is set")

	// Test: name is not set
	var catWithoutName = &SymptomCategory{}
	err = catWithoutName.create(service)
	assert.Error(t, err)

	// Cleanup
	err = service.db.Delete(&cat).Error
	assert.NoError(t, err)
}

func TestCreateSymptomCategoryNoDuplicates(t *testing.T) {
	service := initTest(t)

	var cat = &SymptomCategory{Name: "New Category"}
	var err error
	err = cat.create(service)
	assert.NoError(t, err)

	var duplicateCat = &SymptomCategory{Name: "New Category"}
	err = duplicateCat.create(service)
	assert.NoError(t, err)

	var numberOfCategories int64
	service.db.Find(&SymptomCategory{}, &SymptomCategory{Name: "New Category"}).Count(&numberOfCategories)
	assert.EqualValues(t, 1, numberOfCategories)

	// Cleanup
	err = service.db.Delete(&SymptomCategory{}, &SymptomCategory{Name: "New Category"}).Error
	assert.NoError(t, err)
}

// func TestUpdateSymptomCategory(t *testing.T) {
// 	service, teardown := initTest(t)
// 	defer teardown(t)

// 	service.testCreateConditionEvent(t)

// 	var err error
// 	var cat = &SymptomCategory{ID: 1}
// 	err = cat.update(service, "new name")
// 	assert.NoError(t, err)
// 	var newCat = SymptomCategory{ID: 1}
// 	service.db.Find(&newCat)
// 	assert.Equal(t, "new name", newCat.Name)

// 	// Missing ID
// 	var catMissingID = new(SymptomCategory)
// 	err = catMissingID.update(service, "new name")
// 	assert.Error(t, err)

// 	// Not found
// 	var catNotFound = &SymptomCategory{ID: 1000}
// 	err = catNotFound.update(service, "new name")
// 	assert.Error(t, err)
// }

// func TestDeleteSymptomCategory(t *testing.T) {
// 	service, teardown := initTest(t)
// 	defer teardown(t)

// 	var err error
// 	var cat *SymptomCategory = newSymptomCategory("new name")
// 	service.db.Create(cat)
// 	err = cat.delete(service)
// 	assert.NoError(t, err)

// 	// Cannot delete again: not found
// 	err = cat.delete(service)
// 	assert.Error(t, err)

// 	// Cannot delete with attached symptoms
// 	service.testCreateConditionEvent(t)
// 	cat = &SymptomCategory{ID: 1}
// 	err = cat.delete(service)
// 	assert.Error(t, err)
// }
