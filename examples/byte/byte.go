package main

import (
	"bytes"
	"fmt"

	"github.com/hedhyw/jsoncjson"
)

func main() {
	var in = []byte(`
	{
		/* "If you cannot do great things,
			do small things in a great way."
		                    - Napoleon Hill
		*/
		"Hello": "world"
	}
	`)

	var r = jsoncjson.NewReader(bytes.NewReader(in))

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	fmt.Printf("%s", buf.String())
	// It prints: { "Hello": "world" }
}
