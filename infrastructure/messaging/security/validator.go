package security

import (
	"errors"
	"flowforge-api/infrastructure/config"
)

type WorkerSecurityValidator struct {
	env *config.Config
}

func NewWorkerSecurityValidator(env *config.Config) *WorkerSecurityValidator {
	return &WorkerSecurityValidator{env: env}
}

func (s *WorkerSecurityValidator) Validate(token string) error {
	if token != s.env.RabbitMQSecretKey {
		return errors.New("invalid secret key")
	}
	return nil
}
