package httphelpers

import (
	"fmt"
)

type APIError struct {
	ErrorMessage         error
	ClientErrorMessage   string
	InternalErrorMessage string
}

func NewAPIError(err error, clientMessage string) *APIError {
	return &APIError{
		ErrorMessage:       err,
		ClientErrorMessage: clientMessage,
	}
}

func (e *APIError) SetInternalErrorMessage(msg string) *APIError {
	e.InternalErrorMessage = msg
	return e
}

func (e *APIError) String() string {
	return fmt.Sprintf("%s | %s", e.InternalErrorMessage, e.ErrorMessage)
}