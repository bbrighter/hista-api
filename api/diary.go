package api

import (
	"context"

	entity "encore.app/entity"
)

type DiaryResp struct {
	Diaries []entity.RawDiary `json:"diaries"`
}

// encore:api auth method=GET path=/diary
func (service *Service) GetDiary(ctx context.Context) (DiaryResp, error) {
	meals, events, cats, notes, pollens := service.diary.Get(ctx)
	diaries := entity.CreateRawDiary(meals, events, cats, notes, pollens)
	return DiaryResp{Diaries: diaries}, nil
}
