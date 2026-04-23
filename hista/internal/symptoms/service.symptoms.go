package symptoms

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type SymptomService struct {
	s *SymptomRepo
	c *ConditionRepo
}

func NewSymptomService(db *gorm.DB) *SymptomService {
	s := NewSymptomRepo(db)
	c := NewConditionRepo(db)
	return &SymptomService{s: s, c: c}
}

func (uc SymptomService) ListConditionEvents(ctx context.Context) (ConditionEvents, error) {
	return uc.c.ListConditionEvents(ctx)

}
func (uc SymptomService) CreateConditionEvent(ctx context.Context) (ConditionEvent, error) {
	var event = ConditionEvent{Date: time.Now()}
	err := uc.c.CreateConditionEvent(ctx, &event)
	return event, err

}
func (uc SymptomService) GetConditionEvent(ctx context.Context, id uint) (ConditionEvent, error) {
	return uc.c.FirstConditionEventAndConditions(ctx, id)
}

func (uc SymptomService) PatchConditionEvent(ctx context.Context, id uint, date time.Time) error {
	var values = map[string]any{"date": date}
	return uc.c.UpdateConditionEvent(ctx, id, values)

}
func (uc SymptomService) DeleteConditionEvent(ctx context.Context, id uint) (cats SymptomCategories, err error) {
	err = uc.c.DeleteConditionEvent(ctx, id)
	if err != nil {
		return SymptomCategories{}, err
	}
	cats, err = uc.s.ListSymptomCategoriesAndSymptoms(ctx)
	return cats, err
}

func (uc SymptomService) List(ctx context.Context, eventId uint) ([]Condition, error) {
	conds, err := uc.c.ListConditions(ctx, eventId)
	return conds, err
}

func (uc SymptomService) CreateConditionById(ctx context.Context, eventId uint, symptomId uint) (uint, error) {
	var condition = Condition{ConditionEventID: eventId, SymptomID: symptomId}
	err := uc.c.CreateCondition(ctx, &condition)
	return condition.ID, err
}

func (uc SymptomService) CreateConditionWithNewIngredient(ctx context.Context, eventId uint, symptomName string, symptomCategoryId uint) (uint, SymptomCategories, error) {
	var condition = Condition{
		ConditionEventID: eventId,
		Symptom: Symptom{
			Name:              symptomName,
			SymptomCategoryID: symptomCategoryId,
		}}
	err := uc.c.CreateCondition(ctx, &condition)
	if err != nil {
		return 0, SymptomCategories{}, err
	}
	cats, err := uc.s.ListSymptomCategoriesAndSymptoms(ctx)
	return condition.ID, cats, err

}

func (uc SymptomService) DeleteCondition(ctx context.Context, id uint) (SymptomCategories, error) {
	err := uc.c.DeleteCondition(ctx, id)
	if err != nil {
		return SymptomCategories{}, err
	}
	cats, err := uc.s.ListSymptomCategoriesAndSymptoms(ctx)
	return cats, err
}

func (uc SymptomService) PatchSeverity(ctx context.Context, id uint, newSeverity Severity) error {
	values := map[string]any{"severity": newSeverity}
	return uc.c.UpdateCondition(ctx, id, values)
}

func (uc SymptomService) ListSymptomCategories(ctx context.Context) (SymptomCategories, error) {
	return uc.s.ListSymptomCategoriesAndSymptoms(ctx)
}

func (uc SymptomService) CreateCategory(ctx context.Context, name string) (uint, error) {
	var cat = SymptomCategory{Name: name}
	err := uc.s.CreateSymptomCategory(ctx, &cat)
	return cat.ID, err
}

func (uc SymptomService) ChangeCategory(ctx context.Context, symptomId uint, newCategoryId uint) error {
	values := map[string]any{"symptom_category_id": newCategoryId}
	return uc.s.UpdateSymptom(ctx, symptomId, values)
}
func (uc SymptomService) RenameSymptom(ctx context.Context, symptomId uint, newName string) error {
	values := map[string]any{"name": newName}
	return uc.s.UpdateSymptom(ctx, symptomId, values)
}
func (uc SymptomService) RenameCategory(ctx context.Context, catId uint, newName string) error {
	values := map[string]any{"name": newName}
	return uc.s.UpdateSymptomCategory(ctx, catId, values)
}
func (uc SymptomService) DeleteCategory(ctx context.Context, catId uint) error {
	return uc.s.DeleteSymptomCategory(ctx, catId)
}
