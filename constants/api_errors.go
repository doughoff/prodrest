package constants

import "fmt"

type InvalidParamsError struct {
	Message string
}

func (e *InvalidParamsError) Error() string {
	return e.Message
}

func InvalidParams(message string) error {
	return &InvalidParamsError{
		Message: message,
	}
}

type InvalidBodyError struct {
	Message string
}

func (u InvalidBodyError) Error() string {
	return u.Message
}

func InvalidBody() *InvalidBodyError {
	return &InvalidBodyError{
		fmt.Sprintf("invalid json on request body"),
	}
}
