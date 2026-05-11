package types

import (
	"database/sql/driver"
	"encoding/json"
)

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
