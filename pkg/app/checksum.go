package app

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
)

// Checksum in BigBlueButton authentication system represents an action name, all parameters and a secret concatenated in a single string that is hashed by SHA1.
type Checksum struct {
	Secret string
	Action string
	Params url.Values
}

// StringToSHA1 returns the string value hashed with SHA1 algorithm
func StringToSHA1(value string) (string, error) {
	hasher := sha1.New()

	if _, err := hasher.Write([]byte(value)); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// Value compute the checksum string. It does not hash the value into SHA1 string
func (c *Checksum) Value() string {
	result := c.Action

	for key, element := range c.Params {
		for _, value := range element {
			result += fmt.Sprintf("%s=%s", key, value)
		}
	}
	result += c.Secret

	return result

}
