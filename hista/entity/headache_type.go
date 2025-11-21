package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type HeadacheType string

const (
	Pulsating HeadacheType = "pulsating-pounding"
	Dull      HeadacheType = "dull-pressing"
	Stabbing  HeadacheType = "stabbing"
)

var ValidHeadacheTypes = map[HeadacheType]struct{}{
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
		if _, ok := ValidHeadacheTypes[t]; !ok {
			return fmt.Errorf("invalid headache type: %s", t)
		}
	}
	*e = types
	return nil
}

func (e HeadacheTypes) Value() (driver.Value, error) {
	for _, t := range e {
		if _, ok := ValidHeadacheTypes[t]; !ok {
			return nil, fmt.Errorf("invalid headache type: %s", t)
		}
	}
	return json.Marshal(e)
}
