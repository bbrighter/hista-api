package symptoms

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testEvent *ConditionEvent
var testCondition *Condition
var testSymptom *Symptom
var testCategory *SymptomCategory

func initAPITest(t *testing.T) (*Service, context.Context) {
	var ctx context.Context = context.TODO()
	service := initTest(t)
	return service, ctx
}

func initTest(t *testing.T) *Service {
	service, err := initService()
	assert.NoError(t, err)

	service.initData()

	return service
}

func (service *Service) initData() {
	var event = ConditionEvent{Date: time.Date(2019, 1, 1, 1, 0, 0, 0, time.Local)}
	service.db.FirstOrCreate(&event, &event)
	var symptomCategory = SymptomCategory{Name: "category"}
	service.db.FirstOrCreate(&symptomCategory, &symptomCategory)
	var symptom = Symptom{Name: "symptom", SymptomCategoryID: symptomCategory.ID}
	service.db.FirstOrCreate(&symptom, &symptom)
	var condition = Condition{SymptomID: symptom.ID, ConditionEventID: event.ID, Severity: High}
	service.db.FirstOrCreate(&condition, &condition)

	testEvent = &event
	testCategory = &symptomCategory
	testCondition = &condition
	testSymptom = &symptom
}
