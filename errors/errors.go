package errors

import (
	"bytes"
	"log"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Op describes an operation, usually as the package and method,
// such as "key/server.Lookup".
type Op string

// Separator is the string used to separate nested errors. By
// default, to make errors easier on the eye, nested errors are
// indented on a new line. A server may instead choose to keep each
// error on a single line by modifying the separator string, perhaps
// to ":: ".
var Separator = ":\n\t"

// Kind defines the kind of error this is, mostly for use by systems
// such as FUSE that must act differently depending on the error.
type Kind uint8

const (
	// Unexpected Unclassified error. This value is not printed in the error message.
	Unexpected Kind = iota
	// NotFound Entity or item was not found
	NotFound
	// Forbidden when users inputs bad data
	Forbidden
	// Permission Actor does not have permission to the operation
	Permission
)

func (k Kind) String() string {
	switch k {
	case Unexpected:
		return "unexpected"
	case NotFound:
		return "not found"
	case Permission:
		return "permission denied"
	}

	return "unknown error kind"
}

// Error is the type that implements the error interface.
// It contains a number of fields, each of different type.
// An Error value may leave some values unset.
type Error struct {
	Op       Op
	Kind     Kind
	Err      error
	Severity logrus.Level
}

// Ops prints a "stack trace" of ops
func Ops(e *Error) []Op {
	res := []Op{e.Op}

	subErr, ok := e.Err.(*Error)

	if !ok {
		return res
	}

	res = append(res, Ops(subErr)...)

	return res
}

// E wraps an error
func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}

	e := &Error{}

	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case error:
			e.Err = arg
		case Kind:
			e.Kind = arg
		case logrus.Level:
			e.Severity = arg
		case string:
			e.Err = Str(arg)
		default:
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			panic("bad call to E")
		}
	}

	return e
}

// Str returns an error that formats as the given text. It is intended to
// be used as the error-typed argument to the E function.
func Str(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func (e *Error) isZero() bool {
	return e.Op == "" && e.Kind == 0 && e.Err == nil
}

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}

	b.WriteString(str)
}

func (e *Error) Error() string {
	b := new(bytes.Buffer)

	if e.Op != "" {
		pad(b, ": ")
		b.WriteString(string(e.Op))
	}

	if e.Kind != 0 {
		pad(b, ": ")
		b.WriteString(e.Kind.String())
	}

	if e.Err != nil {
		// Indent on new line if we are cascading non-empty errors.
		if prevErr, ok := e.Err.(*Error); ok {
			if !prevErr.isZero() {
				pad(b, Separator)
				b.WriteString(e.Err.Error())
			}
		} else {
			pad(b, ": ")
			b.WriteString(e.Err.Error())
		}
	}

	if b.Len() == 0 {
		return "no error"
	}

	return b.String()
}

// Is reports whether err is an *Error of the given Kind.
// If err is nil then Is returns false.
func Is(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != Unexpected {
		return e.Kind == kind
	}
	if e.Err != nil {
		return Is(kind, e.Err)
	}
	return false
}
