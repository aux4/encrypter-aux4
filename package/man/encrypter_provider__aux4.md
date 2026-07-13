#### Description

The `aux4` provider command is the local AES-256-GCM encrypter registered under the `aux4/encrypter` facade. It routes to the `encrypt`, `decrypt`, and `generate-secret` implementations.

You normally reach it indirectly: `aux4 encrypter encrypt/decrypt/generate-secret` dispatch here because `aux4` is the default provider. This explicit path is useful for pinning the provider or for debugging.

#### Usage

```bash
aux4 encrypter provider aux4 <command> [options]
```

encrypt          Encrypt plaintext with AES-256-GCM
decrypt          Decrypt base64 ciphertext with AES-256-GCM
generate-secret  Generate a random secret

#### Example

```bash
echo -n 'hello world' | aux4 encrypter provider aux4 encrypt --secret 9j48xBeBCQeV7R6Gnhe1V3acoBSS7ile
```

```text
AQuMlJaRifk/...base64...==
```
