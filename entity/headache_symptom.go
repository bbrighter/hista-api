package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type HeadacheSymptom string

const (
	ShortTermMemory   HeadacheSymptom = "short-term memory"
	Tinnitus          HeadacheSymptom = "tinnitus"
	LightSensitive    HeadacheSymptom = "light-sensitive"
	NoiseSensitive    HeadacheSymptom = "noise-sensitive"
	OdorSensitive     HeadacheSymptom = "odor-sensitive"
	Dizziness         HeadacheSymptom = "dizziness"
	ConcentrationLack HeadacheSymptom = "lack of concentration"
	Tired             HeadacheSymptom = "tired"
	Exhausted         HeadacheSymptom = "exhausted"
	PhysicalActivity  HeadacheSymptom = "physical activity"
	MindActivity      HeadacheSymptom = "mind activity"
)

var validHeadacheSymptoms = map[HeadacheSymptom]struct{}{
	ShortTermMemory: {}, Tinnitus: {}, LightSensitive: {},
	NoiseSensitive: {}, OdorSensitive: {}, Dizziness: {},
	ConcentrationLack: {}, Tired: {}, Exhausted: {},
	PhysicalActivity: {}, MindActivity: {},
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
	*e = symptoms
	return nil
}

func (e HeadacheSymptoms) Value() (driver.Value, error) {
	for _, s := range e {
		if _, ok := validHeadacheSymptoms[s]; !ok {
			return nil, fmt.Errorf("invalid headache symptom: %s", s)
		}
	}
	return json.Marshal(e)
}
