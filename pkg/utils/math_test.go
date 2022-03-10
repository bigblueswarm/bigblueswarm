package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound2Digits(t *testing.T) {
	assert.Equal(t, 16.64, Round2Digits(16.643))
	assert.Equal(t, 16.65, Round2Digits(16.645))
	assert.Equal(t, 16.65, Round2Digits(16.65))
}
