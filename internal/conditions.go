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
	var condition entity.Condition
	var conditionId uint
	var err error
	var symptoms = entity.SymptomCategories{}
	if symptomId != nil {
		conditionId, err = uc.Conditions.CreateConditionBySymptomID(eventId, *symptomId)
	} else if symptomName != nil && symptomCategoryId != nil {
		conditionId, err = uc.Conditions.CreateConditionBySymptomName(eventId, *symptomName, *symptomCategoryId)
		symptoms = uc.Symptoms.ListCategories()
	} else {
		err = errors.ErrorAttributeMustBeSet("symptomId or symptomName and symptomCategoryId")
	}
	if err == nil {
		condition, _ = uc.Conditions.GetCondition(conditionId)
	}
	return condition, symptoms, err
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
