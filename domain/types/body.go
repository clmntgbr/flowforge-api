package types

import (
	"database/sql/driver"
	"encoding"
	"encoding/json"
)

var _ encoding.TextUnmarshaler = (*Body)(nil)

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

func (b *Body) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*b = Body("[]")
		return nil
	}

	return json.Unmarshal(text, (*json.RawMessage)(b))
}

func (b Body) Value() (driver.Value, error) {
	if len(b) == 0 {
		return "[]", nil
	}
	return string(b), nil
}

func (b *Body) Scan(value interface{}) error {
	if value == nil {
		*b = Body("[]")
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

	*b = Body(bytes)
	return nil
}
