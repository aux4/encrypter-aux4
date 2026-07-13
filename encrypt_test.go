package main

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	cases := []string{
		"",
		"hello world",
		"0123456789abcdef",                 // exactly one AES block
		"0123456789abcdef0123456789abcdef", // exactly two AES blocks
		"unicode: café ☕ \n multiline",
	}
	secret := "9j48xBeBCQeV7R6Gnhe1V3acoBSS7ile"

	for _, plaintext := range cases {
		ct, err := encrypt(secret, plaintext)
		if err != nil {
			t.Fatalf("encrypt(%q) error: %v", plaintext, err)
		}
		if ct == "" {
			t.Fatalf("encrypt(%q) produced empty ciphertext", plaintext)
		}
		got, err := decrypt(secret, ct)
		if err != nil {
			t.Fatalf("decrypt error for %q: %v", plaintext, err)
		}
		if got != plaintext {
			t.Fatalf("round-trip mismatch: got %q want %q", got, plaintext)
		}
	}
}

func TestBlockAlignedNonEmpty(t *testing.T) {
	// The old CBC impl produced empty output for 16-byte-aligned input.
	ct, err := encrypt("key", "0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}
	if ct == "" {
		t.Fatal("block-aligned input produced empty ciphertext")
	}
}

func TestVersionByteAndNonceEmbedded(t *testing.T) {
	ct, err := encrypt("key", "x")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := base64.StdEncoding.DecodeString(ct)
	if err != nil {
		t.Fatal(err)
	}
	if raw[0] != versionByte {
		t.Fatalf("expected version byte %d, got %d", versionByte, raw[0])
	}
	// version(1) + nonce(12) + tag(16) minimum for empty-ish plaintext
	if len(raw) < 1+nonceSize+16 {
		t.Fatalf("blob too short: %d", len(raw))
	}
}

func TestWrongKeyFails(t *testing.T) {
	ct, _ := encrypt("rightkey", "secret data")
	_, err := decrypt("wrongkey", ct)
	if err == nil {
		t.Fatal("expected error with wrong key, got nil")
	}
	if !strings.Contains(err.Error(), "wrong key or corrupt") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCorruptInputFails(t *testing.T) {
	_, err := decrypt("key", "not-valid-base64-!!!")
	if err == nil {
		t.Fatal("expected error on corrupt input")
	}
}

func TestTruncatedBlobFails(t *testing.T) {
	// valid base64 but too short to contain version+nonce
	short := base64.StdEncoding.EncodeToString([]byte{0x01, 0x02})
	_, err := decrypt("key", short)
	if err == nil {
		t.Fatal("expected error on truncated blob")
	}
}

func TestKeyDerivationAnyLength(t *testing.T) {
	// SHA-256 derivation means any secret length works (raw AES requires 16/24/32).
	for _, secret := range []string{"a", "short", strings.Repeat("z", 100)} {
		ct, err := encrypt(secret, "data")
		if err != nil {
			t.Fatalf("encrypt with secret len %d: %v", len(secret), err)
		}
		got, err := decrypt(secret, ct)
		if err != nil || got != "data" {
			t.Fatalf("round-trip failed for secret len %d: %v", len(secret), err)
		}
	}
}

func TestGenerateSecretLength(t *testing.T) {
	s, err := generateSecret(24)
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 24 {
		t.Fatalf("expected length 24, got %d", len(s))
	}
}
