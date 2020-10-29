package httphelpers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
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

func (e *APIError) LogError() {
	log.WithFields(log.Fields{
		"Message": e.ErrorMessage,
		"ClientError": e.ClientErrorMessage,
		"InternalErrorMessage": e.InternalErrorMessage,
	}).Error()
}