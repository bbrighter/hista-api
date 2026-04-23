package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/internal/diary"
	"encore.dev/types/uuid"
)

type DiaryResp struct {
	Diaries []diary.Diary `json:"diaries"`
}

// encore:api auth method=GET path=/piid/:piid/diary
func (service *Service) GetDiary(ctx context.Context, piid uuid.UUID) (DiaryResp, error) {
	diaries, err := service.diary.CreateDiary(ctx)
	if err != nil {
		return DiaryResp{}, errors.MapError(err)
	}
	return DiaryResp{Diaries: diaries}, nil
}
