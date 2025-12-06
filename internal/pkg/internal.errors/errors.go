package internalerrors

import (
	"fmt"
	"maps"
	"net/http"
	"strings"
)

type ErrorType string

const (
	ErrorNotFound            ErrorType = "ERR_NOT_FOUND"
	ErrorBadRequest          ErrorType = "ERR_BAD_REQUEST"
	ErrorInternalServerError ErrorType = "ERR_INTERNAL_SERVER"
)

type ErrorDetails map[string]any

type Error struct {
	Type      ErrorType
	Message   string
	Details   ErrorDetails
	RequestID *string
	err       error
}

func New(t ErrorType, reason, msg string) *Error {
	return &Error{
		Type:    t,
		Message: msg,
		Details: ErrorDetails{},
	}
}

func (e *Error) Error() string {
	parts := []string{string(e.Type)}

	if e.Message != "" {
		parts = append(parts, "msg="+e.Message)
	}

	if len(e.Details) > 0 {
		parts = append(parts, fmt.Sprintf("details=%v", e.Details))
	}

	if e.RequestID != nil {
		parts = append(parts, "req="+*e.RequestID)
	}

	if e.err != nil {
		parts = append(parts, "err="+e.err.Error())
	}

	return strings.Join(parts, " | ")
}

func (e *Error) Wrap(err error) *Error {
	e.err = err
	return e
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) WithDetails(m ErrorDetails) *Error {
	if e.Details == nil {
		e.Details = ErrorDetails{}
	}
	maps.Copy(e.Details, m)
	return e
}

func (e *Error) StatusCode() int {
	switch e.Type {
	case ErrorNotFound:
		return http.StatusNotFound
	case ErrorBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func NewInternalServerError(err error) *Error {
	return &Error{
		Type:    ErrorInternalServerError,
		Message: err.Error(),
		Details: ErrorDetails{},
		err:     err,
	}
}
