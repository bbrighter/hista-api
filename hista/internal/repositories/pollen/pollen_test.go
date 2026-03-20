package pollen

import (
	"context"
	"time"

	"encore.app/errors"
	"encore.app/hista/entity"
	"gorm.io/gorm"
)

func (s *PollenRepoTestSuite) TestFindPollenWithSeverity() {
	s.createEvent(time.Now())

	var events entity.PollenEvents
	events = s.repo.FindPollenWithSeverity(0)
	s.Len(events, 1)

	expectedResults := map[entity.PollenType]entity.PollenIntensity{
		entity.Ambrosia: entity.HighPollen,
		entity.Beifuss:  entity.NoPollen,
	}
	results := make(map[entity.PollenType]entity.PollenIntensity)
	for _, p := range events[0].Pollens {
		results[p.Type] = p.Intensity
	}
	s.Equal(expectedResults, results)
}

func (s *PollenRepoTestSuite) TestCreateAndVerify() {
	var err error
	var inputPollens = entity.Pollens{
		{Type: entity.Ambrosia, Intensity: entity.HighPollen},
		{Type: entity.Beifuss, Intensity: entity.MediumToHighPollen},
		{Type: entity.Birke, Intensity: entity.MediumPollen},
		{Type: entity.Esche, Intensity: entity.NoPollen},
		{Type: entity.Graeser, Intensity: entity.NoToSmallPollen},
		{Type: entity.Roggen, Intensity: entity.SmallPollen},
		{Type: entity.Hasel, Intensity: entity.SmallToMediumPollen},
	}
	err = s.repo.Create(s.ctx, inputPollens)
	s.NoError(err)

	pollens, err := gorm.G[entity.Pollen](s.tx).Find(context.Background())
	s.Len(pollens, 7)
	var expectedResult = map[entity.PollenType]entity.PollenIntensity{
		entity.Ambrosia: entity.HighPollen,
		entity.Beifuss:  entity.MediumToHighPollen,
		entity.Birke:    entity.MediumPollen,
		entity.Esche:    entity.NoPollen,
		entity.Graeser:  entity.NoToSmallPollen,
		entity.Roggen:   entity.SmallPollen,
		entity.Hasel:    entity.SmallToMediumPollen,
	}

	var results = make(map[entity.PollenType]entity.PollenIntensity)
	for _, p := range pollens {
		results[p.Type] = p.Intensity
	}

	s.Equal(expectedResult, results)
}

func (s *PollenRepoTestSuite) TestDoesExistAfter() {
	t := time.Date(2017, 11, 3, 4, 0, 0, 0, time.UTC)
	s.createEvent(t)

	err := s.repo.DoesExistAfter(s.ctx, t.Add(time.Minute))
	s.NoError(err)

	err = s.repo.DoesExistAfter(s.ctx, t.Add(-time.Minute))
	s.Error(err)
	s.ErrorIs(err, errors.ErrorAlreadyExists)
}
