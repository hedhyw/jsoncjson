# JSONcJSON

[![Build Status](https://travis-ci.org/hedhyw/jsoncjson.svg?branch=master)](https://travis-ci.org/hedhyw/jsoncjson)

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
