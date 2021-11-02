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
	params := &url.Values{}
	params.Set("name", "supername")
	secret := "supersecret"
	action := "getmeetings"

	checksum := &Checksum{
		Secret: secret,
		Action: action,
		Params: *params,
	}

	assert.Equal(t, checksum.Value(), "getmeetingsname=supernamesupersecret")

}
