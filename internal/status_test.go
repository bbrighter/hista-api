package internal_test

// type TestStatusRepo struct{}

// func (repo TestStatusRepo) Find() entity.Statuses {
// 	return entity.Statuses{}
// }

// func (repo TestStatusRepo) Save(status *entity.Status) error {
// 	status.ID = 1
// 	return nil
// }

// func (repo TestStatusRepo) First(status *entity.Status) error {
// 	status.ID = 1
// 	status.Morning.ID = 2
// 	return nil
// }

// func (repo TestStatusRepo) Delete(status *entity.Status) error {
// 	status.ID = 0
// 	return nil
// }

// func (repo TestStatusRepo) FindForDate(date time.Time) (entity.Status, bool) {
// 	status := entity.Status{
// 		ID:      1,
// 		Date:    time.Now(),
// 		Morning: &entity.MorningStatus{ID: 2},
// 		Evening: &entity.EveningStatus{ID: 3},
// 	}
// 	return status, true
// }

// func TestUpdateMorning(t *testing.T) {
// 	repo := TestStatusRepo{}
// 	uc := internal.NewStatusUseCase(repo)

// 	status, err := uc.SaveMorning(1, time.Now(), 2, entity.Bad, entity.Good)
// 	assert.NoError(t, err)
// 	assert.Equal(t, uint(1), status.ID)
// 	assert.Equal(t, uint(2), status.Morning.ID)
// }
