# JSONcJSON

[![Build Status](https://travis-ci.org/hedhyw/jsoncjson.svg?branch=master)](https://travis-ci.org/hedhyw/jsoncjson)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedhyw/jsoncjson)](https://goreportcard.com/report/github.com/hedhyw/jsoncjson)
[![Coverage Status](https://coveralls.io/repos/github/hedhyw/jsoncjson/badge.svg?branch=master)](https://coveralls.io/github/hedhyw/jsoncjson?branch=master)
[![GoDdoc](https://godoc.org/github.com/hedhyw/jsoncjson?status.svg)](https://godoc.org/github.com/hedhyw/jsoncjson)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/hedhyw/jsoncjson)


The library provides a JSONC (json with comments) to JSON streamer.

It translates JSON with comments:
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

## Usage example:

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
_, = json.NewDecoder(r).Decode(&data)

fmt.Printf("%+v\n", data) // map[Hello:world].
```
