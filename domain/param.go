package domain

import (
	"database/sql/driver"
	"encoding/json"
)

type Param struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Query []Param

type Body json.RawMessage

func (b Body) Value() (driver.Value, error) {
	if b == nil {
		return []byte("null"), nil
	}
	return []byte(b), nil
}

func (b *Body) Scan(value interface{}) error {
	if value == nil {
		*b = Body("null")
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	*b = Body(bytes)
	return nil
}

func (b Body) MarshalJSON() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}
	return []byte(b), nil
}

func (b *Body) UnmarshalJSON(data []byte) error {
	*b = Body(data)
	return nil
}

func (q Query) Value() (driver.Value, error) {
	if q == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(q)
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

type Header []Param

func (h Header) Value() (driver.Value, error) {
	if h == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(h)
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
