package statistics

import (
	"context"

	"encore.app/api/meals"
	"encore.app/notes"
	"encore.app/pollen"
	"encore.app/symptoms"
)

type DiaryResp struct {
	Diaries []RawDiary `json:"diaries"`
}

// encore:api auth method=GET path=/diary
func (service *Service) GetDiary(ctx context.Context) (DiaryResp, error) {
	meals, err := meals.GetMealsAndDependencies(service.db)
	if err != nil {
		return DiaryResp{}, err
	}
	events, cats, err := symptoms.GetConditionEventsAndDependencies(service.db)
	if err != nil {
		return DiaryResp{}, err
	}
	notes := notes.FindAllNotes(service.db)
	pollens := pollen.FindPollenWithSeverity(service.db)
	var input = Input{
		Meals:      meals,
		Events:     events,
		Categories: cats,
		Notes:      notes,
		Pollens:    pollens,
	}
	diaries := createRawDiary(input)
	return DiaryResp{Diaries: diaries}, nil
}
