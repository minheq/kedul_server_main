package errors

import (
	"bytes"
	"fmt"
)

// Kind is category of the error
type Kind int

// Application error categories.
const (
	KindInvalid      Kind = iota + 1 // action cannot be performed or bad request
	KindUnauthorized                 // authorization error
	KindNotFound                     // not found error
	KindUnexpected                   // unexpected error
)

func (kind Kind) String() string {
	switch kind {
	case KindInvalid:
		return "invalid"
	case KindUnauthorized:
		return "unauthorized"
	case KindNotFound:
		return "not found"
	case KindUnexpected:
		return "unexpected"
	}

	return "unknown error kind"
}

// Error defines a standard application error.
type Error struct {
	// Category of the error
	Kind Kind

	// Human-readable message.
	Message string

	// Logical operation and nested error.
	Op string

	// Wrapped error value
	Err error
}

// Invalid returns Error with KindInvalid
func Invalid(op string, err error, message string) *Error {
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
func Unexpected(op string, err error, message string) *Error {
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

	if e.Kind != 0 {
		fmt.Fprintf(&b, "%s: ", e.Kind)
	}

	b.WriteString(e.Message)

	return b.String()
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
func Is(kind Kind, err error) bool {
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
