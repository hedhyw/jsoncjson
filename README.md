# JSONcJSON

![Version](https://img.shields.io/github/v/tag/hedhyw/jsoncjson)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedhyw/jsoncjson)](https://goreportcard.com/report/github.com/hedhyw/jsoncjson)
[![Coverage Status](https://coveralls.io/repos/github/hedhyw/jsoncjson/badge.svg?branch=master)](https://coveralls.io/github/hedhyw/jsoncjson?branch=master)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hedhyw/jsoncjson)](https://pkg.go.dev/github.com/hedhyw/jsoncjson?tab=doc)


The library provides a JSONC (json with comments) to JSON streamer. It
supports multiline comments ( ` /* Comment */ `) and one-line comments
( ` // Comment ` ). It processes chunks of 512 bytes in place.

For example, it translates JSON with comments:
```jsonc
{
    /*
        JSONcJSON
        =^._.^= ∫
    */
    "Hello": "world" // In-line comments are also supported.
}
```
to a valid JSON:
```json
{
    "Hello": "world"
}
```

## Installing:

```sh
go get github.com/hedhyw/jsoncjson
```

## CLI:

Install the command:

```sh
go install github.com/hedhyw/jsoncjson/cmd/jsoncjson@latest
```

It reads JSONC and writes plain JSON. Without file arguments (or with `-`)
the input is read from stdin, the result goes to stdout:

```sh
# From stdin.
echo '{ "Hello": "world" /* Comment. */ }' | jsoncjson

# From a file.
jsoncjson example.jsonc

# To a file.
jsoncjson -o config.json config.jsonc

# Concatenate several inputs.
jsoncjson a.jsonc b.jsonc

# Pipe it further.
jsoncjson config.jsonc | jq .
```

Flags:

| Flag       | Description                        |
| ---------- | ---------------------------------- |
| `-o`       | Output file, `-` for stdout.       |
| `-version` | Print the version and exit.        |
| `-h`       | Print the usage and exit.          |

The command exits with a non-zero status and writes the reason to stderr if
the input cannot be read or the output cannot be written.

## Library usage example:

More [examples](./example_test.go).

```go
// Converting jsonc to json and decoding.

const in = `
{
    "Hello": "world"
    /* Perhaps the truth depends on a walk around the lake. */
}
`

// The reader can be anything.
// For example: file, strings.NewReader(), bytes.NewReader(), ...
var r = jsoncjson.NewReader(strings.NewReader(in))

var data map[string]interface{}
_ = json.NewDecoder(r).Decode(&data)

fmt.Printf("%+v\n", data) // map[Hello:world].
```
