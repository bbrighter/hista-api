package statistics

import (
	"context"

	"encore.app/internal/repositories/meals"
	"encore.app/internal/repositories/notes"
	"encore.app/pollen"
)

type DiaryResp struct {
	Diaries []RawDiary `json:"diaries"`
}

// // encore:api auth method=GET path=/diary
func (service *Service) GetDiary(ctx context.Context) (DiaryResp, error) {
	meals, err := meals.GetMealsAndDependencies(service.db)
	if err != nil {
		return DiaryResp{}, err
	}
	// symrepo := symptoms.NewSymtpomsRepo(service.db)
	// events, cats, err := symrepo.GetConditionEventsAndDependencies()
	if err != nil {
		return DiaryResp{}, err
	}
	repo := notes.NewNotesRepository(service.db)
	notes := repo.List()
	// notes := notes.FindAllNotes(service.db)
	pollens := pollen.FindPollenWithSeverity(service.db)
	var input = Input{
		Meals: meals,
		// Events:     events,
		// Categories: cats,
		Notes:   notes,
		Pollens: pollens,
	}
	diaries := createRawDiary(input)
	return DiaryResp{Diaries: diaries}, nil
}
