package models

import "fmt"

type HttpError struct {
	Message    string
	StatusCode int
}

func NewHttpError(message string, statusCode int) HttpError {
	return HttpError{
		Message:    message,
		StatusCode: statusCode,
	}
}

func NewHttpErrorFromError(message string, err error, statusCode int) HttpError {
	return HttpError{
		Message:    fmt.Sprintf("%s: %v", message, err),
		StatusCode: statusCode,
	}
}

func NewEmptyHttpError() HttpError {
	return HttpError{}
}

func IsHttpErrorEmpty(err HttpError) bool {
	return err == HttpError{}
}
