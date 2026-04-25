package headaches

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type HeadacheType string

func (t HeadacheType) ToString() string {
	return string(t)
}

const (
	Pulsating HeadacheType = "pulsating-pounding"
	Dull      HeadacheType = "dull-pressing"
	Stabbing  HeadacheType = "stabbing"
)

var ValidHeadacheTypes = map[HeadacheType]struct{}{
	Pulsating: {}, Dull: {}, Stabbing: {},
}

func IsValidHeadacheType(s HeadacheType) bool {
	_, ok := ValidHeadacheTypes[s]
	return ok
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
