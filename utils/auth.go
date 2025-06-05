package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
)

// EncodeSHA1Password encode password with key
func EncodeSHA1Password(password string, privateKey string) string {
	if privateKey == "" {
		privateKey = "VinhDhuc@1234a"
	}
	h := sha1.New()
	io.WriteString(h, privateKey+password)

	return fmt.Sprintf("%x", h.Sum(nil))
}
