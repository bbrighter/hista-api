package medicines

import (
	"time"

	"encore.dev/types/uuid"
)

type Medicine struct {
	ID         uint      `gorm:"primaryKey"`
	PIID       uuid.UUID `gorm:"index;type:uuid;not null"`
	Name       string    `gorm:"uniqueIndex"`
	SortOrder  int       `gorm:"not null"`
	IsArchived bool      `gorm:"not null"`
	Intakes    []Intake  `gorm:"constraint:OnDelete:CASCADE"`
}

func (m *Medicine) SetPiid(id uuid.UUID) {
	m.PIID = id
}

type Medicines []*Medicine

type Intake struct {
	ID         uint      `gorm:"primaryKey"`
	Date       time.Time `gorm:"not null"`
	Medicine   Medicine
	MedicineID uint
}

type GroupedIntake struct {
	Date       time.Time `gorm:"column:date"`
	MedicineId uint      `gorm:"column:medicine_id"`
	Count      int64     `gorm:"column:count"`
}

type GroupedIntakeList []GroupedIntake
