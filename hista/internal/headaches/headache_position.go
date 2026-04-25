package headaches

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type HeadachePosition string

func (t HeadachePosition) ToString() string {
	return string(t)
}

const (
	Front    HeadachePosition = "front"
	Back     HeadachePosition = "back"
	Top      HeadachePosition = "top"
	Left     HeadachePosition = "left"
	Right    HeadachePosition = "right"
	Neck     HeadachePosition = "neck"
	Ear      HeadachePosition = "ear"
	Temple   HeadachePosition = "temple"
	Side     HeadachePosition = "side"
	Eye      HeadachePosition = "eye"
	FrontTop HeadachePosition = "front top"
)

var ValidHeadachePositions = map[HeadachePosition]struct{}{
	Front: {}, Back: {}, Top: {},
	Left: {}, Right: {}, Neck: {},
	Ear: {}, Temple: {}, Side: {},
	Eye: {}, FrontTop: {},
}

func IsValidHeadachePosition(s HeadachePosition) bool {
	_, ok := ValidHeadachePositions[s]
	return ok
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
		if _, ok := ValidHeadachePositions[p]; !ok {
			return fmt.Errorf("invalid headache position: %s", p)
		}
	}
	*e = positions
	return nil

}

func (e HeadachePositions) Value() (driver.Value, error) {
	for _, p := range e {
		if _, ok := ValidHeadachePositions[p]; !ok {
			return nil, fmt.Errorf("invalid headache position: %s", p)
		}
	}
	return json.Marshal(e)
}
