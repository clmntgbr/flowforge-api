package errors

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidSignature           = errors.New("invalid signature")
	ErrInvalidRequestBody         = errors.New("invalid request body")
	ErrInvalidEventType           = errors.New("invalid event type")
	ErrValidationFailed           = errors.New("validation failed")
	ErrUserNotAuthenticated       = errors.New("user not authenticated")
	ErrOrganizationNotFound       = errors.New("organization not found")
	ErrUnexpectedSigningMethod    = errors.New("unexpected signing method")
	ErrInvalidToken               = errors.New("invalid token")
	ErrUserNotFound               = errors.New("user not found")
	ErrClerkUserNotFound          = errors.New("clerk user not found")
	ErrUserFailedToCreate         = errors.New("user failed to create")
	ErrUserBanned                 = errors.New("user banned")
	ErrMaxOrganizationsReached    = errors.New("max organizations reached")
	ErrInvalidOrganizationID      = errors.New("invalid organization UUID")
	ErrOrganizationFailedToCreate = errors.New("organization failed to create")
	ErrEndpointsNotFound          = errors.New("endpoints not found")
	ErrInvalidEndpointID          = errors.New("invalid endpoint UUID")
	ErrEndpointNotFound           = errors.New("endpoint not found")
	ErrEndpointFailedToUpdate     = errors.New("endpoint failed to update")
	ErrWorkflowsNotFound          = errors.New("workflows not found")
	ErrInvalidWorkflowID          = errors.New("invalid workflow UUID")
	ErrWorkflowNotFound           = errors.New("workflow not found")
	ErrWorkflowFailedToUpdate     = errors.New("workflow failed to update")
	ErrInvalidRequest             = errors.New("invalid request")
	ErrInvalidStepID              = errors.New("invalid step UUID")
	ErrStepNotFound               = errors.New("step not found")
	ErrStepFailedToUpdate         = errors.New("step failed to update")
	ErrInvalidConnexionID         = errors.New("invalid connexion UUID")
	ErrConnexionNotFound          = errors.New("connexion not found")
	ErrConnexionAlreadyExists     = errors.New("connexion already exists")
)

func GetHTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrUserNotFound),
		errors.Is(err, ErrOrganizationNotFound),
		errors.Is(err, ErrEndpointNotFound),
		errors.Is(err, ErrWorkflowNotFound),
		errors.Is(err, ErrStepNotFound),
		errors.Is(err, ErrConnexionNotFound),
		errors.Is(err, ErrEndpointsNotFound),
		errors.Is(err, ErrWorkflowsNotFound):
		return http.StatusNotFound

	case errors.Is(err, ErrInvalidRequestBody),
		errors.Is(err, ErrValidationFailed),
		errors.Is(err, ErrInvalidOrganizationID),
		errors.Is(err, ErrInvalidEndpointID),
		errors.Is(err, ErrInvalidWorkflowID),
		errors.Is(err, ErrInvalidStepID),
		errors.Is(err, ErrInvalidConnexionID),
		errors.Is(err, ErrInvalidRequest):
		return http.StatusBadRequest

	case errors.Is(err, ErrUserNotAuthenticated),
		errors.Is(err, ErrInvalidToken):
		return http.StatusUnauthorized

	case errors.Is(err, ErrUserBanned):
		return http.StatusForbidden

	case errors.Is(err, ErrMaxOrganizationsReached),
		errors.Is(err, ErrConnexionAlreadyExists):
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}
