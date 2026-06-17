package tag

import (
	"encoding"
	"encoding/json"

	"github.com/google/uuid"
)

type TagInput struct {
	ID    uuid.UUID `json:"id" validate:"required,uuid"`
	Name  string    `json:"name" validate:"required,min=2,max=255"`
	Color string    `json:"color" validate:"required,hexcolor"`
}

type TagInputs []TagInput

var (
	_ encoding.TextUnmarshaler = (*TagInputs)(nil)
	_ json.Unmarshaler         = (*TagInputs)(nil)
)

func (t *TagInputs) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = nil
		return nil
	}

	return json.Unmarshal(data, (*[]TagInput)(t))
}

func (t *TagInputs) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*t = nil
		return nil
	}

	return json.Unmarshal(text, (*[]TagInput)(t))
}
