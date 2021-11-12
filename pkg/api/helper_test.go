package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateHealthCheck(t *testing.T) {
	expected := &HealthCheck{
		ReturnCode: "SUCCESS",
		Version:    "2.0",
	}

	assert.Equal(t, CreateHealthCheck(), expected)
}

func TestDefaultChecksumError(t *testing.T) {
	err := DefaultChecksumError()

	expected := &Error{
		Response{
			ReturnCode: "FAILED",
			MessageKey: "checksumError",
			Message:    "You did not pass the checksum security check",
		},
	}

	assert.Equal(t, err, expected)
}

func TestCreateError(t *testing.T) {
	key := "this is the key error"
	message := "this is the message error"
	err := CreateError(key, message)

	expected := &Error{
		Response{
			ReturnCode: "FAILED",
			MessageKey: key,
			Message:    message,
		},
	}

	assert.Equal(t, err, expected)
}
