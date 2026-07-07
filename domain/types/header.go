package types

import (
	"database/sql/driver"
	"encoding"
	"encoding/json"
)

var (
	_ encoding.TextUnmarshaler = (*Header)(nil)
	_ json.Unmarshaler         = (*Header)(nil)
	_ json.Marshaler           = Header(nil)
)

type Header []Param

func (h Header) MarshalJSON() ([]byte, error) {
	if h == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]Param(h))
}

func (h *Header) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*h = []Param{}
		return nil
	}
	return json.Unmarshal(data, (*[]Param)(h))
}

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

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, (*[]Param)(h))
}

func (h *Header) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*h = []Param{}
		return nil
	}

	return json.Unmarshal(text, (*[]Param)(h))
}
