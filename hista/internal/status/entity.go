package status

import (
	"time"

	"encore.dev/types/uuid"
)

type Status struct {
	ID                    uint      `gorm:"primaryKey"`
	PIID                  uuid.UUID `gorm:"type:uuid;index;not null"`
	Date                  time.Time
	MorningFitness        *int
	DayFitness            *int
	EveningFitness        *int
	MorningSleep          *int
	Depressive            *int
	Tense                 *int
	MoodSwings            *int
	Irritable             *int
	LossOfInterest        *int
	ConcentrationProblems *int
	LackOfDrive           *int
	AppetiteChanges       *int
	SleepProblems         *int
	Overwhelmed           *int
	Crash                 bool
}

func (s *Status) SetPiid(id uuid.UUID) {
	s.PIID = id
}

type Statuses []*Status
