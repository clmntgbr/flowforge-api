package types

import (
	"database/sql/driver"
	"encoding/json"
)

type Body json.RawMessage

func (b Body) MarshalJSON() ([]byte, error) {
	if len(b) == 0 {
		return []byte("[]"), nil
	}
	return json.RawMessage(b).MarshalJSON()
}

func (b *Body) UnmarshalJSON(data []byte) error {
	return (*json.RawMessage)(b).UnmarshalJSON(data)
}

func (b Body) Value() (driver.Value, error) {
	if len(b) == 0 {
		return []byte("[]"), nil
	}
	return []byte(b), nil
}

func (b *Body) Scan(value interface{}) error {
	if value == nil {
		*b = Body("[]")
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	*b = Body(bytes)
	return nil
}
