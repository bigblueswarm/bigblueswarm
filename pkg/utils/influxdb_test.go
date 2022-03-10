package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatInstanceFilter(t *testing.T) {
	type Test struct {
		name      string
		parameter []string
		expected  string
	}
	tests := []Test{
		{
			name:      "No parameter should returns empty string",
			parameter: []string{},
			expected:  "",
		},
		{
			name:      "One parameter should returns a valid filter",
			parameter: []string{"http://localhost:8080"},
			expected:  `r["b3lb_host"] == "http://localhost:8080"`,
		},
		{
			name:      "Multiple parameters should returns a valid filter",
			parameter: []string{"http://localhost:8080", "http://localhost:8081"},
			expected:  `r["b3lb_host"] == "http://localhost:8080" or r["b3lb_host"] == "http://localhost:8081"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, FormatInstancesFilter(test.parameter))
		})
	}
}
