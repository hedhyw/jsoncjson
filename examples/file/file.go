package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/hedhyw/jsoncjson"
)

const exampleFile = "example.jsonc"

func main() {
	var f, err = os.Open(exampleFile)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer f.Close()

	var r = jsoncjson.NewReader(f)

	var data map[string]interface{}
	_ = json.NewDecoder(r).Decode(&data)

	log.Printf("%+v\n", data)
	// It prints: map[Hello:world].
}
