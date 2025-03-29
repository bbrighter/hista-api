package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// @enum short-term memory, tinnitus, light-sensitive, noise-sensitive, odor-sensitive, dizziness, lack of concentration, tired, exhausted
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
