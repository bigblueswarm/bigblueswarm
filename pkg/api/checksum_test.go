package api

import (
	"testing"

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
		parameters string
		action     string
		expected   string
	}

	secret := "supersecret"

	tests := []test{
		{
			name:       "Checksum value with 1 parameter should does not contains any &",
			parameters: "name=supername",
			action:     "getmeetings",
			expected:   "getmeetingsname=supername" + secret,
		},
		{
			name:       "Checksum value with 2 parameters should contains some &",
			parameters: "name=supername&meetingID=1",
			action:     "getmeetings",
			expected:   "getmeetingsname=supername&meetingID=1" + secret,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			checksum := &Checksum{
				Secret: secret,
				Action: test.action,
				Params: test.parameters,
			}

			assert.Equal(t, checksum.Value(), test.expected)
		})
	}
}

func TestChecksumSetTenant(t *testing.T) {
	checksum := &Checksum{
		Params: "param=value",
	}

	checksum.SetTenantMetadata("bbb.localhost.com")
	assert.Equal(t, "param=value&meta_bigblueswarm-tenant=bbb.localhost.com", checksum.Params)
}
