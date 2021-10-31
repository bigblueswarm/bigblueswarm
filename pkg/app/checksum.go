package app

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
)

type Checksum struct {
	Secret string
	Action string
	Params url.Values
}

// params := []gin.Param {
// 	{
// 		Key: "azeae",
// 		Value: "zaeea",
// 	},
// }

func StringToSHA1(value string) (string, error) {
	hasher := sha1.New()

	if _, err := hasher.Write([]byte("getmeetings")); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (c *Checksum) Value() string {
	result := c.Action
	fmt.Println(c.Params)
	for key, element := range c.Params {
		for _, value := range element {
			result += fmt.Sprintf("%s=%s", key, value)
		}
	}
	result += c.Secret

	return result

}
