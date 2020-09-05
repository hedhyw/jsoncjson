package jsoncjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hedhyw/jsoncjson"
)

// This example demonstrates decoding of JSON with comments input with
// a standart "encoding/json" decoder to a map[string]interface{}
// structure.
func Example_jsoncDecode() {
	const jsoncInStr = `
		{
			// A "Hello, World!" program generally is a computer program
			// that outputs or displays the message "Hello, World!".
			
			"Hello": "World"
		}
	`

	var r = jsoncjson.NewReader(strings.NewReader(jsoncInStr))

	var data map[string]interface{}
	_ = json.NewDecoder(r).Decode(&data)

	fmt.Printf("%+v", data)
	// Output: map[Hello:World]
}

// This example demonstrates reading JSON with comments from the bytes
// and printing then JSON without comments output.
func Example_jsoncFromBytes() {
	var in = []byte(`
	{/*
		"If you cannot do great things, do small things in a great way."
		
		- Napoleon Hill

	 */ "Hello": "World" }
	`)

	var r = jsoncjson.NewReader(bytes.NewReader(in))

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	fmt.Printf("%s", buf.String())
	// Output: { "Hello": "World" }
}

// This example demonstrates reading JSON with comments from the string
// and printing then JSON without comments output.
func Example_jsoncFromString() {
	var in = strings.ReplaceAll(`
		{
			// Between saying and doing, many a pair of shoes is worn out.
			"Hello": "World"
		}
	`, "\t", "")

	var r = jsoncjson.NewReader(strings.NewReader(in))

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	fmt.Printf("%s", buf.String())
	// Output:
	// {
	// "Hello": "World"
	// }
}

// This example demonstrates reading JSON with comments from the file
// and decoding with a standart "encoding/json" decoder to a
// map[string]interface{} structure.
func Example_jsoncFromFile() {
	const exampleFile = "./example.jsonc"

	var f, err = os.Open(exampleFile)
	if err != nil {
		log.Fatalf("openning: %s: %s", exampleFile, err)
	}
	defer f.Close()

	var r = jsoncjson.NewReader(f)

	var data map[string]interface{}
	err = json.NewDecoder(r).Decode(&data)
	if err != nil {
		log.Fatalf("decoding: %s", err)
	}

	fmt.Printf("%+v\n", data)
	// Output: map[Hello:world]
}
