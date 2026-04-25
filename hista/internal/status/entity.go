package status

import (
	"time"

	"encore.dev/types/uuid"
)

type Status struct {
	ID             uint      `gorm:"primaryKey"`
	PIID           uuid.UUID `gorm:"type:uuid;index;not null"`
	Date           time.Time
	MorningFitness *int
	EveningFitness *int
	MorningSleep   *int
}

func (s *Status) SetPiid(id uuid.UUID) {
	s.PIID = id
}

type Statuses []*Status
