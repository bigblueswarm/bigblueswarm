package api

import (
	"crypto/sha1"
	"encoding/hex"
)

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
	return c.Action + c.Params + c.Secret
}

// Process compute the value and hash the previous value with SHA1 algorithm
func (c *Checksum) Process() (string, error) {
	return StringToSHA1(c.Value())
}
