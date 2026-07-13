# encrypter-aux4 encrypt

## generate-secret

### should generate a 32-char secret by default

```execute
aux4 encrypter provider aux4 generate-secret | tr -d '\n' | wc -c | tr -d ' '
```

```expect
32
```

### should generate a secret of the requested length

```execute
aux4 encrypter provider aux4 generate-secret 16 | tr -d '\n' | wc -c | tr -d ' '
```

```expect
16
```

## round-trip via provider path

### should encrypt and decrypt back to the original text (stdin)

```execute
KEY=9j48xBeBCQeV7R6Gnhe1V3acoBSS7ile
CT=$(printf 'hello world' | aux4 encrypter provider aux4 encrypt --secret "$KEY")
printf '%s' "$CT" | aux4 encrypter provider aux4 decrypt --secret "$KEY"
```

```expect
hello world
```

### should round-trip a positional text argument

```execute
KEY=9j48xBeBCQeV7R6Gnhe1V3acoBSS7ile
CT=$(aux4 encrypter provider aux4 encrypt --secret "$KEY" 'positional plaintext')
aux4 encrypter provider aux4 decrypt --secret "$KEY" "$CT" </dev/null
```

```expect
positional plaintext
```

### should produce non-empty ciphertext for 16-byte-aligned input

```execute
KEY=9j48xBeBCQeV7R6Gnhe1V3acoBSS7ile
CT=$(printf '0123456789abcdef' | aux4 encrypter provider aux4 encrypt --secret "$KEY")
printf '%s' "$CT" | aux4 encrypter provider aux4 decrypt --secret "$KEY"
```

```expect
0123456789abcdef
```

### should round-trip 32-byte-aligned input

```execute
KEY=9j48xBeBCQeV7R6Gnhe1V3acoBSS7ile
CT=$(printf '0123456789abcdef0123456789abcdef' | aux4 encrypter provider aux4 encrypt --secret "$KEY")
printf '%s' "$CT" | aux4 encrypter provider aux4 decrypt --secret "$KEY"
```

```expect
0123456789abcdef0123456789abcdef
```

## error handling

### should fail cleanly on wrong key (no panic)

```execute
KEY=9j48xBeBCQeV7R6Gnhe1V3acoBSS7ile
CT=$(printf 'secret data' | aux4 encrypter provider aux4 encrypt --secret "$KEY")
printf '%s' "$CT" | aux4 encrypter provider aux4 decrypt --secret 'WRONGKEYWRONGKEYWRONGKEYWRONGKEY'
```

```error:partial
error decrypting: wrong key or corrupt data
```

### should fail cleanly on corrupt input (no panic)

```execute
printf 'not-valid-base64-!!!' | aux4 encrypter provider aux4 decrypt --secret 9j48xBeBCQeV7R6Gnhe1V3acoBSS7ile
```

```error:partial
error decoding base64: **
```

## round-trip via the facade (default provider)

### should encrypt/decrypt through the facade with the default aux4 provider

```execute
KEY=9j48xBeBCQeV7R6Gnhe1V3acoBSS7ile
CT=$(printf 'facade round trip' | aux4 encrypter encrypt --secret "$KEY")
printf '%s' "$CT" | aux4 encrypter decrypt --secret "$KEY"
```

```expect
facade round trip
```
