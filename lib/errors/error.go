package errors

import "fmt"

type ServiceError struct {
	Message string
	Code    int
}

func newServiceError(message string, code int) error {
	return &ServiceError{
		Message: message,
		Code:    code,
	}
}

func (err *ServiceError) Error() string {
	return fmt.Sprintf("error: %v, error code: %v", err.Message, err.Code)
}
