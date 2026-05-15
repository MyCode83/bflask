# bflask

`bflask` is a Go CLI for authorized Flask signed session cookie testing. It bruteforces possible Flask `SECRET_KEY` values from a wordlist and verifies signatures with `github.com/octopart/go-itsdangerous`.

Use it only against applications you own or have explicit permission to test, and in CTF/lab environments.

## Install

```bash
go install .
```

From this source tree:

```bash
go build ./...
go build -o bflask .
go build -o bflask ./cmd/bflask
```

## Usage

```bash
bflask crack \
  -c "eyJ1c2VyIjoiYWRtaW4ifQ..." \
  -w rockyou.txt \
  -t 100
```

JSON output:

```bash
bflask crack -c "$COOKIE" -w keys.txt -j
```

With SHA-256 and a timeout:

```bash
bflask crack -c "$COOKIE" -w keys.txt -d sha256 --timeout 2m
```

Write the hit to a file:

```bash
bflask crack -c "$COOKIE" -w keys.txt -o result.json
```

## Configuration

Priority is `CLI > ENV > CONFIG > DEFAULT`.

Example `config.yaml`:

```yaml
cookie: ""
wordlist: "./wordlist.txt"
threads: 50
salt: "cookie-session"
digest: "sha1"
verbose: false
timeout: "0s"
output: ""
json: false
quiet: false
```

Use a config file:

```bash
bflask --config config.yaml crack
```

Environment variables use the `BFLASK_` prefix:

```bash
BFLASK_THREADS=200 BFLASK_DIGEST=sha256 bflask crack -c "$COOKIE" -w keys.txt
```

## How Flask Cookies Work

Flask's default session stores data client-side and signs it with `SECRET_KEY` through itsdangerous. The data is not encrypted. If the signing key is guessed, the payload can be verified and decoded. `bflask` checks candidate keys with URLSafeTimedSerializer-compatible HMAC signing, `cookie-session` salt by default, and SHA-1 by default.

Compressed Flask cookies start with a leading `.`. `bflask` supports decoding these zlib-compressed payloads when a valid key is found.

## Flags

```text
--cookie, -c string      Flask signed session cookie
--wordlist, -w string    path to SECRET_KEY wordlist
--threads, -t int        number of concurrent workers (default 50)
--salt, -s string        itsdangerous signer salt (default "cookie-session")
--digest, -d string      sha1, sha224, sha256, sha384, sha512, md5 (default "sha1")
--verbose, -v            enable verbose logging
--timeout duration   overall timeout, for example 30s or 5m
--output, -o string      write successful result to a file
--json, -j               emit JSON result
--quiet, -q              print only the recovered SECRET_KEY
```

## Project Tree

```text
.
├── cmd/
├── internal/
│   ├── config/
│   ├── logging/
│   └── output/
├── pkg/
│   └── bflask/
├── examples/
│   └── flask/
├── config.yaml
├── go.mod
├── main.go
└── README.md
```

## Sample Flask App

A tiny app for generating a test cookie is in `examples/flask/app.py`:

```bash
cd examples/flask
python -m flask --app app run
```

Visit `http://127.0.0.1:5000/`, copy the `session` cookie, and test it with a wordlist containing `super secret 1`.

## Troubleshooting

`invalid Flask cookie format`: pass only the cookie value, not `session=<value>` or a full `Cookie:` header.

`unsupported digest`: use one of `sha1`, `sha224`, `sha256`, `sha384`, `sha512`, or `md5`.

No hit found: check the salt, digest, wordlist encoding, and whether the target uses Flask's default signed session serializer.
