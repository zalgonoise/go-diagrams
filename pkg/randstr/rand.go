package randstr

import (
	crand "crypto/rand"
)

const charset = "abcdefghijlkmnopqrstuvwxyz"

func String(length int) string {
	b := make([]byte, length)
	if _, err := crand.Read(b); err != nil {
		return ""
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b)
}
