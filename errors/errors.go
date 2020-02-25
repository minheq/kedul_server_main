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
func Invalid(op string, message string) *Error {
	return &Error{Kind: KindInvalid, Op: op, Message: message}
}

// Unauthorized returns Error with KindUnauthorized
func Unauthorized(op string, err error) *Error {
	return &Error{Kind: KindUnauthorized, Err: err, Op: op}
}

// NotFound returns Error with KindNotFound
func NotFound(op string) *Error {
	return &Error{Kind: KindNotFound, Op: op, Message: "not found"}
}

// Unexpected returns Error with KindUnexpected
func Unexpected(op string, err error, message string) *Error {
	return &Error{Kind: KindUnexpected, Op: op, Err: err, Message: message}
}

// Wrap wraps the inner error
func Wrap(op string, err error, message string) *Error {
	return &Error{Op: op, Err: err, Message: message}
}

// Ops prints the stacktrace
func Ops(e *Error) []string {
	res := []string{e.Op}

	subErr, ok := e.Err.(*Error)

	if !ok {
		return res
	}

	res = append(res, Ops(subErr)...)

	return res
}

func (e *Error) Error() string {
	var b bytes.Buffer

	if e.Op != "" {
		fmt.Fprintf(&b, "%s, %s: ", e.Op, e.Message)
	}

	if e.Err != nil {
		_, ok := e.Err.(*Error)

		if !ok {
			fmt.Fprintf(&b, "%s: ", e.Err.Error())
			return b.String()
		}
		b.WriteString(e.Err.Error())

		return b.String()
	}

	if e.Kind != 0 {
		fmt.Fprintf(&b, "kind: [%s]", e.Kind)
	}

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

// ErrorKind extract error kind from error values
func ErrorKind(err error) Kind {
	if err == nil {
		return 0
	}

	if e, ok := err.(*Error); ok && e.Kind != 0 {
		return e.Kind
	} else if ok && e.Err != nil {
		return ErrorKind(e.Err)
	}

	return KindUnexpected
}

// Is compares whether error matches the kind
func Is(kind Kind, err error) bool {
	return kind == ErrorKind(err)
}
