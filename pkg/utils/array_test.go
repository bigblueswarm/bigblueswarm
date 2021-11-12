package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayContainsString(t *testing.T) {
	assert.True(t, ArrayContainsString([]string{"a", "b", "c"}, "b"))
	assert.False(t, ArrayContainsString([]string{"a", "b", "c"}, "d"))
}
