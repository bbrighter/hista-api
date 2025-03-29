package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Headache struct {
	ID        uint
	Date      time.Time
	Severity  Severity
	Types     HeadacheTypes     `gorm:"type:json"`
	Positions HeadachePositions `gorm:"type:json"`
	Symptoms  HeadacheSymptoms  `gorm:"type:json"`
}

type HeadacheType string

const (
	// pulsierend-pochend, dumpf-drückend, stechend
	Pulsating HeadacheType = "pulsating-pounding"
	Dull      HeadacheType = "dull-pressing"
	Stabbing  HeadacheType = "stabbing"
)

var validHeadacheTypes = map[HeadacheType]struct{}{
	Pulsating: {}, Dull: {}, Stabbing: {},
}

type HeadacheTypes []HeadacheType

func (e *HeadacheTypes) Scan(value any) error {
	if value == nil {
		*e = nil
		return nil
	}
	var types []HeadacheType
	err := unmarshal(value, &types)
	if err != nil {
		return err
	}
	for _, t := range types {
		if _, ok := validHeadacheTypes[t]; !ok {
			return fmt.Errorf("invalid headache type: %s", t)
		}
	}
	*e = types
	return nil
}

func (e HeadacheTypes) Value() (driver.Value, error) {
	for _, t := range e {
		if _, ok := validHeadacheTypes[t]; !ok {
			return nil, fmt.Errorf("invalid headache type: %s", t)
		}
	}
	return json.Marshal(e)
}

type HeadacheSymptom string

const (
	// kurzzeitgedächtnis, ohrenpiepsen, licht-, lärm-, geruchsempfindlich,schwindel, konzentrationsstörung, müde, erschöpft
	ShortTermMemory   HeadacheSymptom = "short-term memory"
	Tinnitus          HeadacheSymptom = "tinnitus"
	LightSensitive    HeadacheSymptom = "light-sensitive"
	NoiseSensitive    HeadacheSymptom = "noise-sensitive"
	OdorSensitive     HeadacheSymptom = "odor-sensitive"
	Dizziness         HeadacheSymptom = "dizziness"
	ConcentrationLack HeadacheSymptom = "lack of concentration"
	Tired             HeadacheSymptom = "tired"
	Exhausted         HeadacheSymptom = "exhausted"
)

var validHeadacheSymptoms = map[HeadacheSymptom]struct{}{
	ShortTermMemory: {}, Tinnitus: {}, LightSensitive: {},
	NoiseSensitive: {}, OdorSensitive: {}, Dizziness: {},
	ConcentrationLack: {}, Tired: {}, Exhausted: {},
}

type HeadacheSymptoms []HeadacheSymptom

func (e *HeadacheSymptoms) Scan(value any) error {
	if value == nil {
		*e = nil
		return nil
	}
	var symptoms []HeadacheSymptom
	err := unmarshal(value, &symptoms)
	if err != nil {
		return err
	}
	for _, s := range symptoms {
		if _, ok := validHeadacheSymptoms[s]; !ok {
			return fmt.Errorf("invalid headache symptom: %s", s)
		}
	}
	return json.Unmarshal(value.([]byte), e)
}

func (e HeadacheSymptoms) Value() (driver.Value, error) {
	for _, s := range e {
		if _, ok := validHeadacheSymptoms[s]; !ok {
			return nil, fmt.Errorf("invalid headache symptom: %s", s)
		}
	}
	return json.Marshal(e)
}

type HeadachePosition string

const (
	Front  HeadachePosition = "front"
	Back   HeadachePosition = "back"
	Both   HeadachePosition = "both"
	Left   HeadachePosition = "left"
	Right  HeadachePosition = "right"
	Neck   HeadachePosition = "neck"
	Ear    HeadachePosition = "ear"
	Temple HeadachePosition = "temple"
)

var validHeadachePositions = map[HeadachePosition]struct{}{
	Front: {}, Back: {}, Both: {},
	Left: {}, Right: {}, Neck: {},
	Ear: {}, Temple: {},
}

type HeadachePositions []HeadachePosition

func (e *HeadachePositions) Scan(value any) error {
	if value == nil {
		*e = nil
		return nil
	}
	var positions []HeadachePosition
	err := unmarshal(value, &positions)
	if err != nil {
		return err
	}
	for _, p := range positions {
		if _, ok := validHeadachePositions[p]; !ok {
			return fmt.Errorf("invalid headache position: %s", p)
		}
	}
	*e = positions
	return nil

}

func (e HeadachePositions) Value() (driver.Value, error) {
	for _, p := range e {
		if _, ok := validHeadachePositions[p]; !ok {
			return nil, fmt.Errorf("invalid headache position: %s", p)
		}
	}
	return json.Marshal(e)
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
