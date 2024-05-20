package statistics

import (
	"context"

	"encore.app/meals"
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
	var input = Input{Meals: meals, Events: events, Categories: cats}
	diaries := diaryFrom(input)
	return DiaryResp{Diaries: diaries}, nil
}
