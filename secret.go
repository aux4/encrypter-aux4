package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generateSecret returns a random alphanumeric string of the given length,
// suitable as AES key material for this provider (any length is accepted by
// deriveKey, but 32 is the sensible default).
func generateSecret(length int) (string, error) {
	secret := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := range secret {
		position, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", fmt.Errorf("error generating secret: %s", err.Error())
		}
		secret[i] = charset[position.Int64()]
	}

	return string(secret), nil
}
