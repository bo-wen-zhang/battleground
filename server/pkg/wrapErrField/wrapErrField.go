package wrapErrField

import "fmt"

// Extends the wrapError to transport log key-value pairs
type wrapErrField struct {
	msg    string
	err    error
	fields []interface{}
}

func (e wrapErrField) Error() string {
	return e.msg
}

func (e wrapErrField) Unwrap() error {
	return e.err
}

func Err(reason string, err error, fields []interface{}) error {
	return &wrapErrField{
		msg:    fmt.Sprintf("%s: %v", reason, err),
		err:    err,
		fields: fields,
	}
}

// Unwraps the fields from the outer most to the inner most
func Fields(err error) []interface{} {
	prev, ok := err.(*wrapErrField)
	if !ok {
		return []interface{}{}
	}
	fields := prev.fields
	for {
		next, ok := prev.err.(*wrapErrField)
		if !ok {
			return fields
		}
		fields = append(append([]interface{}{}, next.fields...), fields...)
		prev = next
	}

}
