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
	ErrClerkUserNotFound       = errors.New("clerk user not found")
	ErrUserFailedToCreate      = errors.New("user failed to create")
	ErrUserBanned              = errors.New("user banned")
	ErrMaxProjectsReached      = errors.New("max projects reached")
	ErrProjectFailedToCreate   = errors.New("project failed to create")
)
