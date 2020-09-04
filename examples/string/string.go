package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hedhyw/jsoncjson"
)

func main() {
	const in = `
	{
		// Between saying and doing, many a pair of shoes is worn out.
		"Hello": "world"
	}
	`

	var r = jsoncjson.NewReader(strings.NewReader(in))

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	fmt.Printf("%s", buf.String())
	// It prints: { "Hello": "world" }
}
