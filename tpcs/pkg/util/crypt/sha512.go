package crypt

import (
	"crypto/sha512"
	"fmt"
)

func EncryptBySHA512(str string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(str)))
}
