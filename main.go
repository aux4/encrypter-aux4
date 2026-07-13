package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printError("Invalid command")
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "generate-secret":
		executeGenerateSecret(args)
	case "encrypt":
		executeEncrypt(args)
	case "decrypt":
		executeDecrypt(args)
	default:
		printError("Invalid command")
		os.Exit(1)
	}
}

func executeGenerateSecret(args []string) {
	lengthParam := "32"
	if len(args) > 0 && args[0] != "" {
		lengthParam = args[0]
	}

	length, err := strconv.Atoi(lengthParam)
	if err != nil {
		printError("Invalid length: %s", lengthParam)
		os.Exit(1)
	}

	secret, err := generateSecret(length)
	if err != nil {
		printError("Error generating secret: %s", err.Error())
		os.Exit(1)
	}

	printOutput(secret)
}

func executeEncrypt(args []string) {
	secret, text := readSecretAndText(args)

	// Plaintext comes from the positional arg if present, otherwise stdin.
	// Encrypt reads stdin raw (no trimming) to preserve the exact bytes.
	if text == "" {
		text = readStdInRaw()
	}

	encrypted, err := encrypt(secret, text)
	if err != nil {
		printError("%s", err.Error())
		os.Exit(1)
	}

	printOutput(encrypted)
}

func executeDecrypt(args []string) {
	secret, text := readSecretAndText(args)

	// Ciphertext comes from the positional arg if present, otherwise stdin.
	if text == "" {
		text = readStdInRaw()
	}

	// base64 is whitespace-free; strip any surrounding whitespace/newlines a
	// pipe (e.g. `echo` without -n) may have added.
	text = strings.TrimSpace(text)

	decrypted, err := decrypt(secret, text)
	if err != nil {
		printError("%s", err.Error())
		os.Exit(1)
	}

	printOutput(decrypted)
}

func readSecretAndText(args []string) (string, string) {
	var secret, text string
	if len(args) > 0 {
		secret = args[0]
	}
	if len(args) > 1 {
		text = args[1]
	}
	return secret, text
}

func readStdInRaw() string {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return ""
	}
	return string(data)
}

func printError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func printOutput(value string) {
	fmt.Fprint(os.Stdout, value)
}
