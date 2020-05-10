package httpError

import "net/http"

type HttpError struct {
	StatusCode  int
	Error       error
	Description []byte
}

func NewError(statusCode int, err error, description []byte) *HttpError {
	return &HttpError{
		StatusCode:  statusCode,
		Error:       err,
		Description: description,
	}
}

func BadRequest(err error, description []byte) *HttpError {
	return NewError(http.StatusBadRequest, err, description)
}

func InternalServerError(err error, description []byte) *HttpError {
	return NewError(http.StatusInternalServerError, err, description)
}

func NotFound(err error, description []byte) *HttpError {
	return NewError(http.StatusNotFound, err, description)
}

func Forbidden(err error, description []byte) *HttpError {
	return NewError(http.StatusForbidden, err, description)
}
