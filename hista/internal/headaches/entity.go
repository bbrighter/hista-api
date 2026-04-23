package headaches

import (
	"encoding/json"
	"fmt"
	"time"

	"encore.dev/types/uuid"
)

type Headache struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	PIID        uuid.UUID `gorm:"index;type:uuid;not null"`
	Date        time.Time
	Severity    uint8
	Types       HeadacheTypes     `gorm:"type:json"`
	Positions   HeadachePositions `gorm:"type:json"`
	Symptoms    HeadacheSymptoms  `gorm:"type:json"`
	Description string
}

func (h *Headache) GetPiid() uuid.UUID {
	return h.PIID
}

func (h *Headache) SetPiid(id uuid.UUID) {
	h.PIID = id
}

type Headaches []*Headache

func unmarshal(value any, output any) error {
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return fmt.Errorf("unexpected type: %T", value)
	}
	if err := json.Unmarshal(bytes, &output); err != nil {
		return err
	}
	return nil
}
