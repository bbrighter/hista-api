package symptoms

import (
	"time"

	"encore.dev/types/uuid"
)

type ConditionEvent struct {
	ID         uint      `gorm:"primaryKey"`
	PIID       uuid.UUID `gorm:"type:uuid;index;not null"`
	Date       time.Time
	Conditions []Condition
}

func (c *ConditionEvent) SetPiid(id uuid.UUID) {
	c.PIID = id
}

type Condition struct {
	ID               uint `gorm:"primaryKey"`
	Severity         Severity
	Symptom          Symptom
	SymptomID        uint
	ConditionEventID uint
}

type Symptom struct {
	ID                uint   `gorm:"primaryKey"`
	Name              string `gorm:"uniqueIndex:idx_name_symptom_category_id"`
	SymptomCategory   SymptomCategory
	SymptomCategoryID uint `gorm:"uniqueIndex:idx_name_symptom_category_id"`
}

type SymptomCategory struct {
	ID       uint      `gorm:"primaryKey"`
	PIID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_name_piid"`
	Name     string    `gorm:"uniqueIndex:idx_name_piid"`
	Symptoms []Symptom
}

func (c *SymptomCategory) SetPiid(id uuid.UUID) {
	c.PIID = id
}

type ConditionEvents []*ConditionEvent

type Conditions []*Condition

type Severity uint8

const (
	VeryLowSeverity  Severity = 1
	LowSeverity      Severity = 2
	MediumSeverity   Severity = 3
	HighSeverity     Severity = 4
	VeryHighSeverity Severity = 5
)

type Symptoms []Symptom

type SymptomCategories []SymptomCategory
