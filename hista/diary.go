package hista

import (
	"context"
	"time"

	"encore.app/errors"
	"encore.app/hista/internal/diary"
	"encore.dev/types/uuid"
)

type DiaryResp struct {
	Date     time.Time `json:"date"`
	Type     string    `json:"type"`
	Content  string    `json:"content"`
	Severity string    `json:"severity"`
	Category string    `json:"category"`
}

type DiaryRespList struct {
	Diaries []DiaryResp `json:"diaries"`
}

func toDiaryRespList(diaries []diary.Diary) DiaryRespList {
	var list = []DiaryResp{}
	for _, d := range diaries {
		list = append(list, DiaryResp{
			Date:     d.Date,
			Type:     string(d.Type),
			Content:  d.Content,
			Severity: d.Severity,
			Category: d.Category,
		})
	}
	return DiaryRespList{Diaries: list}
}

// encore:api auth method=GET path=/piid/:piid/diary
func (service *Service) GetDiary(ctx context.Context, piid uuid.UUID) (DiaryRespList, error) {
	diaries, err := service.diary.CreateDiary(ctx)
	if err != nil {
		return DiaryRespList{}, errors.MapError(err)
	}
	return toDiaryRespList(diaries), nil
}
