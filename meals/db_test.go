package meals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func initTest(t *testing.T) (*Service, func(t *testing.T)) {
	service, err := initService()
	assert.NoError(t, err)

	return service, service.teardown
}

func (service Service) teardown(t *testing.T) {
	var models = []interface{}{&Food{}, &Meal{}, &Ingredient{}}
	var err error
	for _, model := range models {
		err = service.db.Where("1=1").Delete(model).Error
		assert.NoError(t, err)
	}
}

func (service Service) testCreateMeal(t *testing.T) Meal {
	var food = Food{
		ID: 1,
		Ingredient: Ingredient{
			ID:   1,
			Name: "Name",
		},
		IngredientID: 1,
	}
	var meal = Meal{
		ID:    1,
		Date:  time.Now(),
		Foods: []Food{food},
	}
	var err error = service.db.Create(&meal).Error
	assert.NoError(t, err)

	return meal
}
