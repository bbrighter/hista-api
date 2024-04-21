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
	var food = Food{Ingredient: testIngredient, Condition: Cooked}

	id, err = service.createMeal([]Food{food}, date)

	assert.EqualValues(t, 1, id)
	assert.NoError(t, err)

	food = Food{Ingredient: Ingredient{Name: "Name", ID: 100000}, Condition: Cooked}
	id, err = service.createMeal([]Food{food}, date)
	assert.Error(t, err)
	// var count int64
	// service.db.Find(&Meal{}).Count(&count)
	// assert.EqualValues(t, 1, count)
	// service.db.Find(&Food{}).Count(&count)
	// assert.EqualValues(t, 1, count)
	// service.db.Find(&Ingredient{}).Count(&count)
	// assert.EqualValues(t, 1, count)
}

func TestGetMeal(t *testing.T) {
	service, _ := initService()

	_, err := service.getMeal(1000)
	assert.Error(t, err)

	// _, err = service.getMeal(1)
	// assert.NoError(t, err)
}
