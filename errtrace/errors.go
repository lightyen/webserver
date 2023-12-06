package errtrace

import (
	"errors"
)

type Error struct {
	err   error
	stack StackTrace
}

func (e *Error) Unwrap() error { return e.err }

func (e *Error) Error() string { return e.err.Error() }

func (e *Error) Stack() string {
	return e.stack.String()
}

func (e *Error) String() string {
	return e.Error() + "\n" + e.Stack()
}

func WrapError(err error) *Error {
	return &Error{err: err, stack: Stack(3, 6)}
}

func IsError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*Error)
	if ok {
		return true
	}
	return IsError(errors.Unwrap(err))
}

func AsError(err error) (*Error, bool) {
	var e *Error
	ok := errors.As(err, &e)
	if !ok {
		return nil, false
	}
	return e, true
}

type RecoveredError struct {
	err   any
	stack StackTrace
}

func (e *RecoveredError) Unwrap() error {
	if v, ok := e.err.(error); ok {
		return v
	}
	return nil
}

func (e *RecoveredError) Error() string {
	switch v := e.err.(type) {
	default:
		return "<unknown error>"
	case string:
		return v
	case error:
		return v.Error()
	}
}

func (e *RecoveredError) String() string {
	return e.Error() + "\n" + e.Stack()
}

func (e *RecoveredError) Stack() string { return e.stack.String() }

// Wrap the return value from recover() func
func WrapRecoveredError(r any) *RecoveredError {
	return &RecoveredError{err: r, stack: Stack(5, 10)}
}
