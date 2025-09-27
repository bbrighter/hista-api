package entity

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"encore.dev/types/uuid"
)

type Headache struct {
	ID          uint
	Date        time.Time
	Severity    HeadacheSeverity
	Types       HeadacheTypes     `gorm:"type:json"`
	Positions   HeadachePositions `gorm:"type:json"`
	Symptoms    HeadacheSymptoms  `gorm:"type:json"`
	Description string
	PIID        uuid.UUID `gorm:"type:uuid;index"`
}

func (h *Headache) GetPiid() uuid.UUID {
	return h.PIID
}

func (h *Headache) SetPiid(id uuid.UUID) {
	h.PIID = id
}

type HeadacheSeverity uint8

type Headaches []*Headache

func (h Headache) ToResp() HeadacheResponse {
	return HeadacheResponse{
		ID:          h.ID,
		Date:        h.Date,
		Severity:    h.Severity,
		Types:       h.Types,
		Positions:   h.Positions,
		Symptoms:    h.Symptoms,
		Description: h.Description,
	}
}

type HeadacheResponse struct {
	ID          uint              `json:"id"`
	Date        time.Time         `json:"date"`
	Severity    HeadacheSeverity  `json:"severity"`
	Types       HeadacheTypes     `json:"types" encore:"optional"`
	Positions   HeadachePositions `json:"positions" encore:"optional"`
	Symptoms    HeadacheSymptoms  `json:"symptoms" encore:"optional"`
	Description string            `json:"description"`
}

type HeadachesResponse struct {
	Headaches []HeadacheResponse `json:"headaches"`
}

func (hs Headaches) ToResp() HeadachesResponse {
	var headaches = []HeadacheResponse{}
	for _, h := range hs {
		headaches = append(headaches, h.ToResp())
	}
	sort.Slice(headaches, func(i, j int) bool {
		return headaches[i].Date.After(headaches[j].Date)
	})
	return HeadachesResponse{Headaches: headaches}
}

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
