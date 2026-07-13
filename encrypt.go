package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// versionByte is an internal, format-version marker prepended inside the
// encrypted blob (before the nonce). It never appears in a human-visible way —
// it is base64-encoded together with the rest of the blob. Bumping it lets a
// future format change be detected on decrypt without breaking the value shape.
const versionByte = 0x01

// nonceSize is the standard AES-GCM nonce length (12 bytes / 96 bits).
const nonceSize = 12

// deriveKey turns arbitrary secret input into a fixed 32-byte AES-256 key.
// The current impl fed the secret bytes directly to aes.NewCipher, which
// required the caller to supply exactly 16/24/32 bytes. Hashing with SHA-256
// is the robust option: any secret length works and always yields a valid
// 256-bit key.
func deriveKey(secret string) []byte {
	sum := sha256.Sum256([]byte(secret))
	return sum[:]
}

// encrypt seals text with AES-256-GCM and returns
// base64( versionByte || nonce || ciphertext||tag ). No padding (GCM is a
// stream mode, so 16-byte-aligned input is not a special case), no external IV.
func encrypt(secret string, text string) (string, error) {
	block, err := aes.NewCipher(deriveKey(secret))
	if err != nil {
		return "", fmt.Errorf("error creating cipher: %s", err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM: %s", err.Error())
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error generating nonce: %s", err.Error())
	}

	// prefix: [versionByte][nonce...] ; Seal appends ciphertext+tag onto it.
	blob := make([]byte, 1+nonceSize)
	blob[0] = versionByte
	copy(blob[1:], nonce)

	sealed := gcm.Seal(blob, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// decrypt reverses encrypt. On a wrong key, corrupt/truncated input, or
// tampering it returns a clean error (never panics) — GCM authentication
// fails and cipher.Open returns an error instead of garbage.
func decrypt(secret string, text string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", fmt.Errorf("error decoding base64: %s", err.Error())
	}

	if len(decoded) < 1+nonceSize {
		return "", fmt.Errorf("error decrypting: input too short or corrupt")
	}

	if decoded[0] != versionByte {
		return "", fmt.Errorf("error decrypting: unsupported format version %d", decoded[0])
	}

	nonce := decoded[1 : 1+nonceSize]
	ciphertext := decoded[1+nonceSize:]

	block, err := aes.NewCipher(deriveKey(secret))
	if err != nil {
		return "", fmt.Errorf("error creating cipher: %s", err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM: %s", err.Error())
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting: wrong key or corrupt data")
	}

	return string(plaintext), nil
}
