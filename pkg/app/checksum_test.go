package app

import (
	"testing"

	"net/url"

	"github.com/stretchr/testify/assert"
)

func TestStringToSHA1(t *testing.T) {
	expected := "7936e787e9ea1fb449c7f767d3a76fe092e63cb0"
	value := "getmeetings"
	sha1, err := StringToSHA1(value)
	if err != nil {
		assert.Fail(t, err.Error())
	} else {
		assert.Equal(t, sha1, expected)
	}
}

func TestChecksumValue(t *testing.T) {
	type test struct {
		name       string
		parameters *url.Values
		action     string
		expected   string
	}

	secret := "supersecret"

	tests := []test{
		{
			name: "Checksum value with 1 parameter should does not contains any &",
			parameters: &url.Values{
				"name": []string{"supername"},
			},
			action:   "getmeetings",
			expected: "getmeetingsname=supername" + secret,
		},
		{
			name: "Checksum value with 2 parameters should contains some &",
			parameters: &url.Values{
				"name":      []string{"supername"},
				"meetingID": []string{"1"},
			},
			action:   "getmeetings",
			expected: "getmeetingsname=supername&meetingID=1" + secret,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			checksum := &Checksum{
				Secret: secret,
				Action: test.action,
				Params: *test.parameters,
			}

			assert.Equal(t, checksum.Value(), test.expected)
		})
	}
}
