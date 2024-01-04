package wrapErrField

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	// zerolog supports including key-value pairs in logs using Fields()
	"github.com/rs/zerolog"
)

var mockError = errors.New("some error")

var mockFields = []interface{}{"name", "John", "age", 23, "gender", "M"}

var mockWrappedError = WrapMsgAndFields("some message", mockError, mockFields)

func getMockWrappedError() error {
	return mockWrappedError
}

func TestFromFuncReturn(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"name\":\"John\",\"age\":23,\"gender\":\"M\"}\n"

	err := getMockWrappedError()

	logger.Info().Fields(Fields(err)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}

func TestErrorf(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"name\":\"John\",\"age\":23,\"gender\":\"M\"}\n"

	err := fmt.Errorf("some error: %w", mockWrappedError)

	logger.Info().Fields(Fields(err)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}

func TestUnwrap(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"name\":\"John\",\"age\":23,\"gender\":\"M\"}\n"

	err := fmt.Errorf("some %s error: %w", "", mockWrappedError)

	err = errors.Unwrap(err)

	logger.Info().Fields(Fields(err)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}

func TestWrapFields(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"name\":\"John\",\"age\":23,\"gender\":\"M\"}\n"

	err := WrapFields(errors.New("some error"), mockFields)

	logger.Info().Fields(Fields(err)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}

func TestWrapFieldsUnwrap(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\"}\n"

	err := WrapFields(errors.New("some error"), mockFields)

	err = errors.Unwrap(err)
	logger.Info().Fields(Fields(err)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}

func TestPropogate(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"name\":\"John\",\"age\":23,\"gender\":\"M\"}\n"

	err := func() error {
		err := func() error {
			return WrapFields(mockError, []interface{}{"age", 23, "gender", "M"})
		}()
		return WrapFields(err, []interface{}{"name", "John"})
	}()

	logger.Info().Fields(Fields(err)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}

func TestUnwrapJoinedErrs(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"name\":\"John\",\"age\":23,\"gender\":\"M\"}\n"

	err1 := WrapFields(errors.New("error 1"), []interface{}{"name", "John", "age", 23})
	err2 := WrapFields(errors.New("error 1"), []interface{}{"gender", "M"})

	err := fmt.Errorf("%w %w", err1, err2)

	logger.Info().Fields(Fields(err)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}
