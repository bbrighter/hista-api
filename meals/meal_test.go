package meals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateMeal(t *testing.T) {
	service, _ := initService()

	var id uint
	var err error

	var testIngredient = Ingredient{ID: 100, Name: "Test ingredient"}
	service.db.Create(&testIngredient)
	var date = time.Now()

	id, err = service.createMeal(date)

	assert.EqualValues(t, 1, id)
	assert.NoError(t, err)

	// Cleanup
	service.deleteMeal(id)
}

func TestGetMeal(t *testing.T) {
	service, _ := initService()

	_, err := service.getMeal(1000)
	assert.Error(t, err)

	id, _ := service.createMeal(time.Now())
	_, err = service.getMeal(id)
	assert.NoError(t, err)

	service.deleteMeal(id)
}

func TestDeleteMeal(t *testing.T) {
	service, _ := initService()

	var err error
	err = service.deleteMeal(1)
	assert.Error(t, err)

	id, _ := service.createMeal(time.Now())
	err = service.deleteMeal(id)
	assert.NoError(t, err)

}
