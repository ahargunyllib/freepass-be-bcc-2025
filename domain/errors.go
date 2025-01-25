package domain

import "errors"

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
	StatusCode: 409,
	Err:        errors.New("email already exists"),
}

var ErrUserNotFound = &RequestError{
	StatusCode: 404,
	Err:        errors.New("user not found"),
}

var ErrCannotDeleteAdmin = &RequestError{
	StatusCode: 403,
	Err:        errors.New("cannot delete admin"),
}
