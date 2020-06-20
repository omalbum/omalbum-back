package messages

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

type Message struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	httpStatus int
}

func (m Message) Error() string {
	return fmt.Sprintf("code:%s, message: %s", m.Code, m.Message)
}

func New(code string, message string) error {
	return Message{
		Code:    code,
		Message: message,
	}
}

func newWithHttpStatus(code string, message string, httpStatus int) error {
	return Message{
		Code:       code,
		Message:    message,
		httpStatus: httpStatus,
	}
}

func NewBadRequest(code string, message string) error {
	return newWithHttpStatus(code, message, http.StatusBadRequest)
}

func NewConflict(code string, message string) error {
	return newWithHttpStatus(code, message, http.StatusConflict)
}

func NewNotFound(code string, message string) error {
	return newWithHttpStatus(code, message, http.StatusNotFound)
}

// Returns the httpCode
func GetHttpCode(err error) int {
	if message, ok := err.(Message); ok {
		return message.httpStatus
	}

	return http.StatusOK
}

// Returns the code
func GetCode(err error) string {
	if message, ok := err.(Message); ok {
		return message.Code
	}

	return ""
}

// Returns the first error from a ozzo validation error in a proper way
func NewValidation(err error) error {
	if errors, okErrors := err.(validation.Errors); okErrors {
		return getError(errors, "")
	}

	return nil
}

// Recursive function
func getError(errors validation.Errors, previousKey string) error {
	for key, obj := range errors {
		if errorObject, okObject := obj.(validation.ErrorObject); okObject {
			var composedCode string

			if previousKey == "" {
				composedCode = key + "_" + errorObject.Code()
			} else {
				composedCode = previousKey + "_" + key + "_" + errorObject.Code()
			}

			return Message{
				Code:    composedCode,
				Message: errorObject.Error(),
			}
		}
		if errorsObject, okObject := obj.(validation.Errors); okObject {
			return getError(errorsObject, key)
		}
	}

	return nil
}
