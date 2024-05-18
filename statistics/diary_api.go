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
	meals := meals.GetMealsAndDependencies(service.db)
	events, cats := symptoms.GetConditionEventsAndDependencies(service.db)
	var input = Input{Meals: meals, Events: events, Categories: cats}
	diaries := diaryFrom(input)
	return DiaryResp{Diaries: diaries}, nil
}
