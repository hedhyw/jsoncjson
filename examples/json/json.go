package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hedhyw/jsoncjson"
)

func main() {
	const in = `
	{
		"Hello": "world"
		/* Perhaps the truth depends on a walk around the lake. */
	}
	`

	var r = jsoncjson.NewReader(strings.NewReader(in))

	var data map[string]interface{}
	_ = json.NewDecoder(r).Decode(&data)

	fmt.Printf("%+v\n", data)
	// It prints: map[Hello:world].
}
