package statistics

import "context"

type DiaryResp struct {
	Diaries []RawDiary `json:"diaries"`
}

// encore:api auth method=GET path=/diary
func (service *Service) GetDiary(ctx context.Context) (DiaryResp, error) {
	meals, events := getDiaryData(service)
	diaries := diaryFrom(meals, events)
	return DiaryResp{Diaries: diaries}, nil
}
