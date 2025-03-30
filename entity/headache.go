package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

type Headache struct {
	ID          uint
	Date        time.Time
	Severity    HeadacheSeverity
	Types       HeadacheTypes     `gorm:"type:json"`
	Positions   HeadachePositions `gorm:"type:json"`
	Symptoms    HeadacheSymptoms  `gorm:"type:json"`
	Description string
}

type HeadacheSeverity uint8

type Headaches []Headache

func (h Headache) ToResp() HeadacheResponse {
	return HeadacheResponse(h)
}

type HeadacheResponse struct {
	ID          uint              `json:"id"`
	Date        time.Time         `json:"date"`
	Severity    HeadacheSeverity  `json:"severity"`
	Types       HeadacheTypes     `json:"types"`
	Positions   HeadachePositions `json:"positions"`
	Symptoms    HeadacheSymptoms  `json:"symptoms"`
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
