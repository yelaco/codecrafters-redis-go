package util

import (
	"crypto/rand"
	"math/big"
)

const alphanumeric = "abcdefghijklmnopqrstuvwxyz123456789"

// RandomAlphanumericString generates a random alphanumeric string of the specified length
func RandomAlphanumericString(length int) (string, error) {
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphanumeric))))
		if err != nil {
			return "", err
		}
		b[i] = alphanumeric[num.Int64()]
	}
	return string(b), nil
}
