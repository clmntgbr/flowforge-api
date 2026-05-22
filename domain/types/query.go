package types

import (
	"database/sql/driver"
	"encoding/json"
)

type Query []Param

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
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, q)
}
