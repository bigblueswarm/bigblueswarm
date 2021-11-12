package utils

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestComputeErr(t *testing.T) {
	type test struct {
		name     string
		error    error
		expected error
	}

	tests := []test{
		{
			name:     "Computing error should returns the error",
			error:    fmt.Errorf("Failed to compute error"),
			expected: fmt.Errorf("Failed to compute error"),
		},
		{
			name:     "Computing error should returns nil value",
			error:    redis.Nil,
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, ComputeErr(test.error), test.expected)
		})
	}
}
