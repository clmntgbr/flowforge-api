package presenter

import "encoding/json"

type ClerkUserCreated struct {
	ID        string `json:"id" validate:"required"`
	FirstName string `json:"first_name" validate:"omitempty"`
	LastName  string `json:"last_name" validate:"omitempty"`
	Banned    *bool  `json:"banned" validate:"required"`
}

type ClerkUser struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Banned    bool   `json:"banned"`
}

type ClerkUserUpdated struct {
	ID        string `json:"id" validate:"required"`
	FirstName string `json:"first_name" validate:"omitempty"`
	LastName  string `json:"last_name" validate:"omitempty"`
	Banned    *bool  `json:"banned" validate:"required"`
}

type ClerkUserDeleted struct {
	ID string `json:"id" validate:"required"`
}

type ClerkEvent struct {
	Type       string          `json:"type" validate:"required"`
	InstanceID string          `json:"instance_id" validate:"required"`
	Object     string          `json:"object" validate:"required"`
	Timestamp  int64           `json:"timestamp" validate:"required"`
	Data       json.RawMessage `json:"data"`
}
