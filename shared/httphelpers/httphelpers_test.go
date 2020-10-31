package httphelpers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIErrorSuccess(t *testing.T) {
	err := errors.New("This is a test error")
	clientErrorMessage := ""

	successAPIError := &APIError{err, clientErrorMessage, ""}

	apiError := NewAPIError(err, clientErrorMessage)
	assert.Equal(t, successAPIError, apiError)
}

func TestSetInternalErrorMessageSuccess(t *testing.T) {
	err := errors.New("This is a test error")
	clientErrorMessage := "Client error message test"
	internalErrorMessage := "This is an internal message test"

	successAPIError := &APIError{err, clientErrorMessage, ""}
	successAPIError.SetInternalErrorMessage(internalErrorMessage)

	assert.Equal(t, successAPIError.InternalErrorMessage, internalErrorMessage)
}
