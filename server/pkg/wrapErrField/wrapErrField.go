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

func (e *wrapError) Error() string {
	return e.msg
}

func (e *wrapError) Unwrap() error {
	return e.err
}

// Fields are unwrapped and flattened from the top of the stack to the bottom.
// This function is capable of unwrapping joined errors that implement Unwrap() []error
// The order of fields return is deterministic
func Fields(err error) []interface{} {
	curr := err
	fields := []interface{}{}
	for {
		if curr == nil {
			return fields
		}
		//switch curr.(type) {
		if _, ok := curr.(*wrapError); ok {
			fields = append(fields, curr.(*wrapError).fields...)
		} else {
			curr, ok := curr.(interface {
				Unwrap() []error
			})
			if ok && curr != nil {
				fmt.Println("Hello")
				for _, e := range curr.Unwrap() {
					fields = append(fields, Fields(e)...)
				}
			}
		}
		curr = errors.Unwrap(curr)
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
