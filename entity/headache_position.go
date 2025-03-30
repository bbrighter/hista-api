package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// @enum front, back, both, left, right, neck, ear, temple
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
