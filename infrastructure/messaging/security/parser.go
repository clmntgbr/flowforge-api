package security

import (
	"bytes"
	"encoding/json"
	"errors"
	"flowforge-api/infrastructure/config"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
)

const MaxWorkerPayloadBytes = 1024 * 1024

type WorkerParser struct {
	env *config.Config
}

func NewWorkerParser(env *config.Config) *WorkerParser {
	return &WorkerParser{env: env}
}

func (s *WorkerParser) ParseAndValidate(body []byte, dest any) error {
	if len(body) == 0 {
		return errors.New("message body is empty")
	}

	if len(body) > MaxWorkerPayloadBytes {
		return fmt.Errorf("message body exceeds max size (%d bytes)", MaxWorkerPayloadBytes)
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dest); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); err != io.EOF {
		return errors.New("invalid JSON: unexpected trailing content")
	}

	validator := validator.New()
	if err := validator.Struct(dest); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	return nil
}
