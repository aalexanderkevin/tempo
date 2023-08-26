package model

import (
	"errors"
	"fmt"
)

const (
	ErrorBadRequest          int = 400
	ErrorUnauthorized        int = 403
	ErrorNotFound            int = 404
	ErrorDuplicate           int = 409
	ErrorUnprocessableEntity int = 422
	ErrorInternalServer      int = 500
)

type Error struct {
	Message string `json:"error_message"`
	Code    int
}

func NewError(message string, code int) Error {
	return Error{
		Message: message,
		Code:    code,
	}
}

func (e Error) Error() string {
	return e.Message
}

func NewParameterError(msg *string) Error {
	defaultMessage := "invalid parameter"
	if msg == nil {
		msg = &defaultMessage
	}
	return NewError(*msg, ErrorUnprocessableEntity)
}

func NewNotFoundError() Error {
	return NewError("resource not found", ErrorNotFound)
}

func NewDuplicateError() Error {
	return NewError("resource already exists", ErrorDuplicate)
}

func NewUnauthorizedError() Error {
	return NewError("unauthorized access", ErrorUnauthorized)
}

func NewBadRequestError(msg *string) Error {
	defaultMessage := "bad request"
	if msg == nil {
		msg = &defaultMessage
	}
	return NewError(*msg, ErrorBadRequest)
}

func IsDuplicateError(e error) bool {
	var internalErr Error
	if !errors.As(e, &internalErr) {
		return false
	}

	return internalErr.Code == ErrorDuplicate
}

func IsNotFoundError(e error) bool {
	var internalErr Error
	if !errors.As(e, &internalErr) {
		return false
	}

	return internalErr.Code == ErrorNotFound
}

func NewStatusNotOKError(code int, body []byte) Error {
	e := fmt.Sprintf("status is not ok, status=%d body=%s", code, body)
	return NewError(e, ErrorInternalServer)
}

func IsParameterError(e error) bool {
	var internalErr Error
	if !errors.As(e, &internalErr) {
		return false
	}

	return internalErr.Code == ErrorUnprocessableEntity
}

func IsBadRequestError(e error) bool {
	var internalErr Error
	if !errors.As(e, &internalErr) {
		return false
	}

	return internalErr.Code == ErrorBadRequest
}
