# aux4 encrypter-aux4

Local **AES-256-GCM** encryption provider for the [`aux4/encrypter`](https://hub.aux4.io/aux4/encrypter) facade. It registers the `aux4` provider, so `aux4 encrypter encrypt/decrypt` work out of the box with no external services.

Installing this package pulls in `aux4/encrypter` automatically via `dependencies`.

## Installation

```bash
aux4 aux4 pkger install aux4/encrypter-aux4
```

## Usage

Through the facade (recommended — `aux4` is the default provider):

```bash
SECRET=$(aux4 encrypter generate-secret)

echo -n 'hello world' | aux4 encrypter encrypt --secret "$SECRET"
# -> base64 ciphertext

echo -n '<base64>' | aux4 encrypter decrypt --secret "$SECRET"
# -> hello world
```

Directly (explicit provider path):

```bash
echo -n 'hello world' | aux4 encrypter provider aux4 encrypt --secret "$SECRET"
aux4 encrypter provider aux4 decrypt --secret "$SECRET" '<base64>'
```

## Cryptography

- **Algorithm:** AES-256-GCM (authenticated encryption).
- **Key derivation:** the `--secret` value is hashed with SHA-256 to a 256-bit key. Any secret length is accepted — unlike raw AES, which requires exactly 16/24/32 bytes. This is the robust choice and means a `generate-secret` of any length works as a key.
- **Nonce:** a fresh random 12-byte nonce is generated per encryption and embedded in the blob. There is **no `--iv`** parameter.
- **No padding:** GCM is a stream mode, so plaintext whose length is a multiple of the AES block size (16 bytes) is not a special case — it encrypts correctly and produces non-empty output. (This fixes a bug in the previous AES-CBC implementation.)
- **Integrity:** GCM's authentication tag means a wrong key, corrupt, or tampered blob is detected. Decryption then exits **non-zero** with a clean error message — never a panic. (This fixes the previous decrypt-panic bug.)

### Value format

The emitted value is plain base64 with **no prefix**:

```text
base64( versionByte(1) ‖ nonce(12) ‖ ciphertext‖tag )
```

The single internal `versionByte` (currently `0x01`) is invisible in the base64 output; it lets a future format change be detected on decrypt without changing the value's shape. There is no `--iv` and no human-visible `aux4enc:` prefix.

## Commands

| Command | Description |
|---------|-------------|
| `encrypt` | Encrypt plaintext (positional arg or stdin) → base64 on stdout |
| `decrypt` | Decrypt base64 (positional arg or stdin) → plaintext on stdout |
| `generate-secret` | Generate a random alphanumeric secret (default length 32) |

These implement the [Provider Command Contract](https://hub.aux4.io/aux4/encrypter) defined by `aux4/encrypter`, under the `encrypter:provider:aux4` profile.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `AUX4_SECRET` | Fallback for `--secret` (the AES key material) |
