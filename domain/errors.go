package domain

import (
	"errors"
	"net/http"
)

type SerializableError interface {
	Serialize() any
}

type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return r.Err.Error()
}

var ErrEmailAlreadyExists = &RequestError{
	StatusCode: http.StatusConflict,
	Err:        errors.New("email already exists"),
}

var ErrUserNotFound = &RequestError{
	StatusCode: http.StatusNotFound,
	Err:        errors.New("user not found"),
}

var ErrCannotDeleteAdmin = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("cannot delete admin"),
}

var ErrInvalidCredentials = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("invalid credentials"),
}

var ErrNoBearerToken = &RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("no bearer token provided"),
}

var ErrInvalidBearerToken = &RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("invalid bearer token"),
}

var ErrExpiredBearerToken = &RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("expired bearer token"),
}

var ErrBearerTokenNotActive = &RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("bearer token not active"),
}

var ErrCantAccessResource = &RequestError{
	StatusCode: http.StatusForbidden,
	Err:        errors.New("you don't have access to this resource"),
}

var ErrSessionNotFound = &RequestError{
	StatusCode: http.StatusNotFound,
	Err:        errors.New("session not found"),
}

var ErrSessionCannotBeUpdated = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("session cannot be updated"),
}

var ErrSessionCannotBeDeleted = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("session cannot be deleted"),
}

var ErrSessionProposalLimit = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("session proposal limit reached"),
}

var ErrSessionAlreadyStarted = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("session already started"),
}

var ErrSessionNotStarted = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("session not started"),
}

var ErrSessionAlreadyEnded = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("session already ended"),
}

var ErrSessionNotEnded = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("session not ended"),
}

var ErrSessionNotAccepted = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("session not accepted"),
}

var ErrSessionFull = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("session is full"),
}

var ErrSessionAlreadyRegistered = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("session already registered"),
}

var ErrSessionNotRegistered = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("session not registered"),
}

var ErrSessionTimeConflict = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("session time conflict"),
}

var ErrSessionCancelled = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("already cancel this session"),
}

var ErrReviewDeleted = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("review already deleted"),
}

var ErrSessionAlreadyReviewed = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("session already reviewed"),
}

var ErrClaimsNotFound = &RequestError{
	StatusCode: http.StatusInternalServerError,
	Err:        errors.New("claims not found"),
}

var ErrCantUpdateTitle = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("you can't update the title"),
}
