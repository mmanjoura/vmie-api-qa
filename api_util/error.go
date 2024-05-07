// This package provides various utility functions for working with interfaces, maps, 
// arrays, and JSON/YAML conversion. Each function has a comment explaining its 
// purpose and usage. Let me know if you need further clarification on any specific function!

package api_util


import (
	"fmt"
	"runtime/debug"
)

const (
	ErrOK         = iota // 0
	ErrInvalid           // invalid parameters
	ErrNotFound          // resource not found
	ErrExpect            // the REST result doesn't match the expected value
	ErrHttp              // Http request failed
	ErrServerResp        // unexpected server response
	ErrInternal          // unexpected internal error (meqa error)
)

// Error implements MQ specific error type.
type Error interface {
	error
	Type() int
}

// TypedError holds a type and a back trace for easy debugging
type TypedError struct {
	errType int
	errMsg  string
}

// Error returns the error message associated with the TypedError.
func (e *TypedError) Error() string {
	return e.errMsg
}

// Type returns the error type of the TypedError.
func (e *TypedError) Type() int {
	return e.errType
}

// NewError creates a new error with the specified error type and error message.
// It also includes a backtrace of the error stack.
// The error type is an integer that represents the type of the error.
// The error message is a string that describes the error in detail.
// The function returns an error interface.
func NewError(errType int, str string) error {
	buf := string(debug.Stack())
	err := TypedError{errType, ""}
	err.errMsg = fmt.Sprintf("==== %v ====\nError message:\n%s\nBacktrace:%v", errType, str, buf)
	return &err
}