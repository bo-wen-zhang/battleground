package wrapErrField

import (
	"errors"
	"strings"
	"testing"

	// zerolog supports including key-value pairs in logs using Fields()
	"github.com/rs/zerolog"
)

var ErrSomeError = errors.New("some error")

func TestManyFields(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"name\":\"John\",\"age\":23,\"gender\":\"M\"}\n"

	err := func() error {
		fields := []interface{}{
			"name", "John",
			"age", 23,
			"gender", "M",
		}
		return WrapMsgAndFields("some message", ErrSomeError, fields)
	}()

	logger.Info().Fields(Fields(err)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}

func TestMultipleFields(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"name\":\"John\",\"age\":23}\n"

	err := func() error {
		fields := []interface{}{
			"name", "John", "age", 23,
		}
		return WrapMsgAndFields("some message", ErrSomeError, fields)
	}()

	logger.Info().Fields(Fields(err)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}

func TestCallStackHeightTwo(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"age\":23,\"name\":\"John\"}\n"

	err := func() error {
		return WrapMsgAndFields("some message", ErrSomeError, []interface{}{"name", "John"})
	}()
	newErr := WrapMsgAndFields("some message", err, []interface{}{"age", 23})

	logger.Info().Fields(Fields(newErr)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}

func TestUnwrap(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"name\":\"John\"}\n"

	err := func() error {
		return WrapMsgAndFields("some message", ErrSomeError, []interface{}{"name", "John"})
	}()
	newErr := WrapMsgAndFields("some message", err, []interface{}{"age", 23})

	newErr = Unwrap(newErr)
	logger.Info().Fields(Fields(newErr)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}

func TestPropogate(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"name\":\"John\"}\n"

	err := func() error {
		err := WrapFields(func() error {
			return WrapMsgAndFields("some message", ErrSomeError, []interface{}{"name", "John"})
		}(), []interface{}{})
		return err
	}()

	logger.Info().Fields(Fields(err)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}
