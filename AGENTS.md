# AGENTS.md

Guidance for AI coding agents working in this repository.

## What this project is

`jsoncjson` (module `github.com/hedhyw/jsoncjson`) is a tiny, dependency-free
Go library that streams JSONC (JSON with comments) into valid JSON. It exposes
an `io.Reader` wrapper that strips one-line comments (`// ...`) and multiline
comments (`/* ... */`) on the fly, so the output can be fed directly into
`encoding/json` or any other consumer. Input is processed in chunks of
`bytes.MinRead` (512) bytes; comments inside JSON strings are preserved, and
escape sequences are handled correctly.

## Public API (import path `github.com/hedhyw/jsoncjson`)

The entire public API is a single function:

- `jsoncjson.NewReader(r io.Reader) io.Reader` — wraps any reader that yields
  JSONC and returns a reader that yields pure JSON.

Everything else (`jsoncTranslator`, the token state machine) is unexported.
Keep it that way — new exported surface needs a strong justification.

## Layout

```
jsoncjson.go        # the whole implementation: NewReader + token state machine
jsoncjson_test.go   # unit tests (table-driven, stdlib testing only)
example_test.go     # runnable godoc examples (Example_*)
example.jsonc       # fixture used by Example_jsoncFromFile
makefile            # lint/test targets
```

Single root package, no subpackages, no dependencies (`go.mod` has no
`require` block).

## Commands

```sh
make lint           # golangci-lint (installed into ./bin on first run)
make test           # go test
make test.coverage  # go test with coverage profile (coverage.out)
make                # lint + test + install git pre-commit hook
```

Go version: see `go.mod`. There is no code generation.

## Conventions

- Zero dependencies: do not add any `require` to `go.mod`.
- The reader is streaming and allocation-light — avoid buffering the whole
  input; work byte-by-byte through the existing state machine
  (`handleToken`) in `jsoncjson.go`.
- Tests use only the standard library (no assertion frameworks). New behavior
  needs a test in `jsoncjson_test.go`; user-visible usage patterns belong in
  `example_test.go` as `Example_*` functions with `// Output:` blocks.
- CI (`.github/workflows/check.yml`) runs `make lint` and `make test.coverage`
  on pushes and pull requests; PR titles must follow Conventional Commits
  (checked by `.github/workflows/semantic.yaml`).
