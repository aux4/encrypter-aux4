#### Description

The `decrypt` command decrypts a base64-encoded AES-256-GCM blob and prints the recovered plaintext to stdout.

- The `--secret` value is hashed to a 256-bit AES key with SHA-256 (must match the key used to encrypt).
- The nonce is read from inside the blob — there is no `--iv`.
- The ciphertext is taken from the positional `text` argument if present, otherwise from stdin. Surrounding whitespace is ignored.
- On a wrong key, corrupt input, or tampering, GCM authentication fails: the command exits non-zero with a clean error message on stderr and never panics.

#### Usage

```bash
aux4 encrypter provider aux4 decrypt --secret <keyInput> [text]
```

--secret  AES key material (any length; hashed to 256 bits). Env: AUX4_SECRET
text      The base64 ciphertext to decrypt (positional arg; falls back to stdin)

#### Example

```bash
echo -n 'AQuMlJaRifk/...base64...==' | aux4 encrypter provider aux4 decrypt --secret 9j48xBeBCQeV7R6Gnhe1V3acoBSS7ile
```

```text
hello world
```
