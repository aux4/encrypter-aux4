#### Description

The `encrypt` command encrypts plaintext with AES-256-GCM and prints a base64-encoded blob to stdout.

- The `--secret` value is hashed to a 256-bit AES key with SHA-256, so any secret length is accepted.
- A fresh random 12-byte nonce is generated per call and embedded in the output. There is no `--iv`.
- The plaintext is taken from the positional `text` argument if present, otherwise from stdin.
- Output is plain base64 (`versionByte ‖ nonce ‖ ciphertext‖tag`) with no prefix. GCM is a stream mode, so block-aligned input is not a special case.

#### Usage

```bash
aux4 encrypter provider aux4 encrypt --secret <keyInput> [text]
```

--secret  AES key material (any length; hashed to 256 bits). Env: AUX4_SECRET
text      The plaintext to encrypt (positional arg; falls back to stdin)

#### Example

```bash
echo -n 'hello world' | aux4 encrypter provider aux4 encrypt --secret 9j48xBeBCQeV7R6Gnhe1V3acoBSS7ile
```

```text
AQuMlJaRifk/...base64...==
```
