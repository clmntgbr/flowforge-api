package errors

import "errors"

var (
	ErrInvalidSignature        = errors.New("invalid signature")
	ErrInvalidRequestBody      = errors.New("invalid request body")
	ErrInvalidEventType        = errors.New("invalid event type")
	ErrValidationFailed        = errors.New("validation failed")
	ErrUserNotAuthenticated    = errors.New("user not authenticated")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
	ErrUserNotFound            = errors.New("user not found")
	ErrUserBanned              = errors.New("user banned")
)
