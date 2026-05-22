package types

import (
	"database/sql/driver"
	"encoding/json"
)

type Header []Param

func (h Header) Value() (driver.Value, error) {
	if h == nil {
		return "[]", nil
	}
	bytes, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (h *Header) Scan(value interface{}) error {
	if value == nil {
		*h = []Param{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, h)
}
