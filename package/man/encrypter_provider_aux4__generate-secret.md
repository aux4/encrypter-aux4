#### Description

The `generate-secret` command generates a random alphanumeric secret and prints it to stdout. It is suitable as `--secret` key material for the `aux4` provider (any length works, since the secret is hashed to a 256-bit key).

#### Usage

```bash
aux4 encrypter provider aux4 generate-secret [length]
```

length  The length of the secret (positional arg; default: 32)

#### Example

```bash
aux4 encrypter provider aux4 generate-secret 32
```

```text
Mh2SxkdnEScloyqj8VzEgzJlBgE3kfT9
```
