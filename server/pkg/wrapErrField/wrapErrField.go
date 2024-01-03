// I need to add comments to this package :)
package wrapErrField

import (
	"errors"
	"fmt"
)

// Enables transporting key-value pairs up a call stack for logging purposes.
type wrapError struct {
	msg    string
	err    error
	fields []interface{}
}

func (e wrapError) Error() string {
	return e.msg
}

func Unwrap(err error) error {
	if prev, ok := err.(*wrapError); ok {
		return prev.err
	}
	return errors.Unwrap(err)
}

// Fields are unwrapped from the top of the stack to the bottom.
func Fields(err error) []interface{} {
	prev, ok := err.(*wrapError)
	if !ok {
		return []interface{}{}
	}
	fields := prev.fields
	for {
		next, ok := prev.err.(*wrapError)
		if !ok {
			return fields
		}
		fields = append(fields, next.fields...)
		prev = next
	}
}

func WrapFields(err error, fields []interface{}) error {
	return &wrapError{
		msg:    err.Error(),
		err:    err,
		fields: fields,
	}
}

func WrapMsgAndFields(msg string, err error, fields []interface{}) error {
	return &wrapError{
		msg:    fmt.Sprintf("%s: %v", msg, err),
		err:    err,
		fields: fields,
	}
}
