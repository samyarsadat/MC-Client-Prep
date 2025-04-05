package testutils

import (
	"crypto/rand"
	"encoding/base64"
)

func GetRandomBase64(length uint64) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)
}
