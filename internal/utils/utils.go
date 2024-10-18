package utils

import (
	"crypto/sha512"
	"encoding/base64"
)

func Sha512Sum(input string) string {
	hasher := sha512.New()
	hasher.Write([]byte(input))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}
