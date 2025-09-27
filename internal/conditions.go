package internal

import (
	"context"

	"encore.app/entity"
)

type (
	IConditionRepo interface {
		ListConditions(ctx context.Context, eventId uint) ([]*entity.Condition, error)
		CreateConditionBySymptomName(ctx context.Context, eventId uint, symptomName string, categoryId uint) (uint, error)
		CreateConditionBySymptomID(ctx context.Context, eventId uint, symptomId uint) (uint, error)
		DeleteCondition(ctx context.Context, conditionId uint) error
		ChangeSeverity(ctx context.Context, conditionId uint, newSeverity entity.Severity) error
		GetCondition(ctx context.Context, id uint) (entity.Condition, error)
	}

	IConditionUseCase interface {
		List(ctx context.Context, eventId uint) (entity.Conditions, error)
		Create(ctx context.Context, eventId uint, symptomName *string, symptomId *uint, symptomCategoryId *uint) (entity.Condition, entity.SymptomCategories, error)
		Delete(ctx context.Context, id uint) (entity.SymptomCategories, error)
		PatchSeverity(ctx context.Context, id uint, newSeverity entity.Severity) error
	}
)

type ConditionUseCase struct {
	Conditions IConditionRepo
	Symptoms   ISymptomCategoriesRepo
}

func NewConditionsUseCase(conditions IConditionRepo, symptoms ISymptomCategoriesRepo) ConditionUseCase {
	return ConditionUseCase{
		Conditions: conditions,
		Symptoms:   symptoms,
	}
}

func (uc ConditionUseCase) List(ctx context.Context, eventId uint) (entity.Conditions, error) {
	conds, err := uc.Conditions.ListConditions(ctx, eventId)
	return conds, errorMapper(err)
}

func (uc ConditionUseCase) Create(ctx context.Context, eventId uint, symptomName *string, symptomId *uint, symptomCategoryId *uint) (entity.Condition, entity.SymptomCategories, error) {
	var err error
	var conditionId uint
	if symptomId != nil {
		conditionId, err = uc.Conditions.CreateConditionBySymptomID(ctx, eventId, *symptomId)
	} else if symptomName != nil && symptomCategoryId != nil {
		conditionId, err = uc.Conditions.CreateConditionBySymptomName(ctx, eventId, *symptomName, *symptomCategoryId)
	}
	if err != nil {
		return entity.Condition{}, entity.SymptomCategories{}, errorMapper(err)
	}

	symptoms, err := uc.Symptoms.ListCategories(ctx)
	if err != nil {
		return entity.Condition{}, entity.SymptomCategories{}, errorMapper(err)
	}
	condition, err := uc.Conditions.GetCondition(ctx, conditionId)
	return condition, symptoms, errorMapper(err)
}

func (uc ConditionUseCase) Delete(ctx context.Context, id uint) (entity.SymptomCategories, error) {
	err := uc.Conditions.DeleteCondition(ctx, id)
	if err != nil {
		return entity.SymptomCategories{}, errorMapper(err)
	}
	cats, err := uc.Symptoms.ListCategories(ctx)
	return cats, errorMapper(err)
}

func (uc ConditionUseCase) PatchSeverity(ctx context.Context, id uint, newSeverity entity.Severity) error {
	return errorMapper(uc.Conditions.ChangeSeverity(ctx, id, newSeverity))
}
