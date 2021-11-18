package errors

import (
	"bytes"
	"fmt"
)

type Error struct {
	// The kind of error
	Kind Kind `json:"kind"`
	// The method or function being invoked
	Method Method `json:"method"`
	// The username of the user attempting the operation
	Username Username `json:"username"`
	// The error message
	Message Message `json:"message"`
	// Nested Error
	Err error `json:"err"`
}

func E(args ...interface{}) error {
	if len(args) == 0 {
		return nil
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Kind:
			e.Kind = arg
		case Method:
			e.Method = arg
		case Username:
			e.Username = arg
		case Message:
			e.Message = arg
		case string:
			if e.Message == "" {
				e.Message = Message(arg)
			}
		case *Error:
			copy := *arg
			e.Err = &copy
		case error:
			e.Err = arg
		default:
			return Errorf("unknown type %T with value %v in error call", arg, arg)
		}
	}
	return e
}

func (e *Error) isZero() bool {
	return e.Method == "" && e.Username == "" && e.Message == "" && e.Kind == 0 && e.Err == nil
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
	if e.Method != "" {
		b.WriteString(string(e.Method))
	}
	if e.Username != "" {
		pad(b, ": ")
		b.WriteString(string(e.Username))
	}
	if e.Kind != 0 {
		pad(b, ": ")
		b.WriteString(e.Kind.String())
	}
	if e.Message != "" {
		pad(b, ": ")
		b.WriteString(string(e.Message))
	}
	if e.Err != nil {
		// Indent to new line if it is another Error
		if prevErr, ok := e.Err.(*Error); ok {
			if !prevErr.isZero() {
				pad(b, ":\n\t")
				b.WriteString(string(e.Err.Error()))
			}
		} else { // Just print it out if a standard error
			pad(b, ": ")
			b.WriteString(string(e.Err.Error()))
		}
	}
	if b.Len() == 0 {
		return "no error message"
	}
	return b.String()
}

// Errorf is a wrapper that calls E with just a message formatted with the args
// this allows for only importing this library for all error handling
func Errorf(f string, args ...interface{}) error {
	return E(fmt.Sprintf(f, args...))
}

// Returns True if err is of type Error and of the given Kind
func Is(k Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != Other {
		return e.Kind == k
	}
	if e.Err != nil {
		return Is(k, e.Err)
	}
	return false
}
