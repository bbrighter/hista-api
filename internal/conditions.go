package internal

import (
	"encore.app/entity"
	"encore.app/errors"
)

type ConditionUseCase struct {
	Conditions IConditionRepo
	Symptoms   ISymptomCategoriesRepo
}

func NewconditionsUseCase(conditions IConditionRepo, symptoms ISymptomCategoriesRepo) ConditionUseCase {
	return ConditionUseCase{
		Conditions: conditions,
		Symptoms:   symptoms,
	}
}

func (uc ConditionUseCase) List(eventId uint) entity.Conditions {
	return uc.Conditions.ListConditions(eventId)
}

func (uc ConditionUseCase) Create(eventId uint, symptomName *string, symptomId *uint, symptomCategoryId *uint) (entity.Condition, entity.SymptomCategories, error) {
	var condition = &entity.Condition{
		ConditionEventID: eventId,
		Severity:         entity.MediumSeverity,
	}
	var err error
	var symptoms = entity.SymptomCategories{}
	if symptomId != nil {
		condition.SymptomID = *symptomId
		err = uc.Conditions.CreateConditionBySymptomID(condition)
	} else if symptomName != nil && symptomCategoryId != nil {
		err = uc.Conditions.CreateConditionBySymptomName(condition, *symptomName, *symptomCategoryId)
		symptoms = uc.Symptoms.ListCategories()
	} else {
		err = errors.ErrorAttributeMustBeSet("symptomId or symptomName and symptomCategoryId")
	}
	if err != nil {
		return *condition, symptoms, err
	}
	*condition, err = uc.Conditions.GetCondition(condition.ID)
	return *condition, symptoms, err
}

func (uc ConditionUseCase) Delete(id uint) (entity.SymptomCategories, error) {
	err := uc.Conditions.DeleteCondition(id)
	if err != nil {
		return entity.SymptomCategories{}, err
	}
	return uc.Symptoms.ListCategories(), nil
}

func (uc ConditionUseCase) PatchSeverity(id uint, newSeverity entity.ConditionSeverity) error {
	return uc.Conditions.ChangeSeverity(id, newSeverity)
}
