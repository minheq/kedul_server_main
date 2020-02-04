package errors

import (
	"bytes"
	"fmt"
)

// Application error categories.
const (
	KindInvalid      = "invalid"      // action cannot be performed or bad request
	KindUnauthorized = "unauthorized" // authorization error
	KindUnexpected   = "unexpected"   // unexpected error
	KindNotFound     = "not found"    // unexpected error
)

// Error defines a standard application error.
type Error struct {
	// Category of the error
	Kind string

	// Human-readable message.
	Message string

	// Logical operation and nested error.
	Op string

	// Wrapped error value
	Err error
}

// Invalid returns Error with KindInvalid
func Invalid(op string, message string) *Error {
	return &Error{Kind: KindInvalid, Op: op, Message: message}
}

// Unauthorized returns Error with KindUnauthorized
func Unauthorized(op string) *Error {
	return &Error{Kind: KindUnauthorized, Op: op}
}

// NotFound returns Error with KindNotFound
func NotFound(op string) *Error {
	return &Error{Kind: KindNotFound, Op: op}
}

// Unexpected returns Error with KindUnexpected
func Unexpected(op string, err error) *Error {
	return &Error{Kind: KindUnexpected, Op: op, Err: err}
}

func (e *Error) Error() string {
	var b bytes.Buffer

	if e.Op != "" {
		fmt.Fprintf(&b, "%s: ", e.Op)
	}

	if e.Err != nil {
		b.WriteString(e.Err.Error())

		return b.String()
	}

	if e.Kind != "" {
		fmt.Fprintf(&b, "<%s> ", e.Kind)
	}

	b.WriteString(e.Message)

	return b.String()
}

// ErrorKind returns the code of the root error, if available. Otherwise returns KindUnexpected.
func ErrorKind(err error) string {
	if err == nil {
		return ""
	}

	if e, ok := err.(*Error); ok && e.Kind != "" {
		return e.Kind
	} else if ok && e.Err != nil {
		return ErrorKind(e.Err)
	}

	return KindUnexpected
}

// ErrorMessage extract messages from error values
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	if e, ok := err.(*Error); ok && e.Message != "" {
		return e.Message
	} else if ok && e.Err != nil {
		return ErrorMessage(e.Err)
	}
	return "An internal error has occurred. Please contact technical support."
}

// Is compares whether error matches the kind
func Is(kind string, err error) bool {
	e, ok := err.(*Error)

	if !ok {
		return false
	}

	if e.Kind != KindUnexpected {
		return e.Kind == kind
	}

	if e.Err != nil {
		return Is(kind, e.Err)
	}

	return false
}
