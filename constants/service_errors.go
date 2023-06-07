package constants

import "fmt"

type UniqueConstraintError struct {
	Message string
}

func (u UniqueConstraintError) Error() string {
	return u.Message
}

func NewUniqueConstrainError(field string) *UniqueConstraintError {
	return &UniqueConstraintError{
		fmt.Sprintf("unique constraint error on: %v", field),
	}
}

type RequiredFieldError struct {
	Message string
}

func (u RequiredFieldError) Error() string {
	return u.Message
}

func NewRequiredFieldError(field string) *RequiredFieldError {
	return &RequiredFieldError{
		fmt.Sprintf("required field: %v", field),
	}
}

type InvalidOperationError struct {
	Message string
}

func (u InvalidOperationError) Error() string {
	return u.Message
}

func NewInvalidOperationError(additionalInfo string) *InvalidOperationError {
	return &InvalidOperationError{
		fmt.Sprintf("required field: %v", additionalInfo),
	}
}
