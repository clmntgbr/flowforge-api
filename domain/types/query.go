package types

import (
	"database/sql/driver"
	"encoding"
	"encoding/json"
)

var (
	_ encoding.TextUnmarshaler = (*Query)(nil)
	_ json.Unmarshaler         = (*Query)(nil)
	_ json.Marshaler           = Query(nil)
)

type Query []Param

func (q Query) MarshalJSON() ([]byte, error) {
	if q == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]Param(q))
}

func (q *Query) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*q = []Param{}
		return nil
	}
	return json.Unmarshal(data, (*[]Param)(q))
}

func (q Query) Value() (driver.Value, error) {
	if q == nil {
		return "[]", nil
	}
	bytes, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (q *Query) Scan(value interface{}) error {
	if value == nil {
		*q = []Param{}
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

	return json.Unmarshal(bytes, (*[]Param)(q))
}

func (q *Query) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*q = []Param{}
		return nil
	}

	return json.Unmarshal(text, (*[]Param)(q))
}
