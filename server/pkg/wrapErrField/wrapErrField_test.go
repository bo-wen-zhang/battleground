package wrapErrField

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

var ErrSomeError = errors.New("some error")

func TestSingleStrField(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"name\":\"John\"}\n"

	err := func() error {
		fields := []interface{}{
			"name", "John",
		}
		return Err("some message", ErrSomeError, fields)
	}()

	logger.Info().Fields(Fields(err)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}

func TestSingleIntField(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"age\":23}\n"

	err := func() error {
		fields := []interface{}{
			"age", 23,
		}
		return Err("some message", ErrSomeError, fields)
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
		return Err("some message", ErrSomeError, fields)
	}()

	logger.Info().Fields(Fields(err)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}

func TestCallStackHeightTwo(t *testing.T) {
	b := strings.Builder{}
	logger := zerolog.New(&b)
	expected := "{\"level\":\"info\",\"name\":\"John\",\"age\":23}\n"

	err := func() error {
		fields := []interface{}{"name", "John"}
		return Err("some message", ErrSomeError, fields)
	}()

	newFields := []interface{}{"age", 23}
	newErr := Err("some message", err, newFields)
	fmt.Println(err, newErr)
	logger.Info().Fields(Fields(newErr)).Send()

	if b.String() != expected {
		t.Errorf("Logged %s, expected %s", b.String(), expected)
	}
}
