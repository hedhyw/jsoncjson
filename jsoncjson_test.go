package jsoncjson_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/hedhyw/jsoncjson"
)

func TestJSONWithoutComments(t *testing.T) {
	const jsoncInStr = `{"hello": "world"}`
	var jsoncIn = strings.NewReader(jsoncInStr)

	var jsoncOut = jsoncjson.NewReader(jsoncIn)
	var buff bytes.Buffer
	var n, err = buff.ReadFrom(jsoncOut)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var jsonOutStr = buff.String()
	if int64(len(jsonOutStr)) != n {
		t.Fatalf("exp: %d, got: %d", len(jsonOutStr), n)
	}

	if jsoncInStr != jsonOutStr {
		t.Fatalf("exp: %s, got: %s", jsoncInStr, jsonOutStr)
	}
}

func TestJSONWithCommentsReader(t *testing.T) {
	type testData struct {
		Hello string `json:"hello"`
	}

	var exp = testData{
		Hello: "World",
	}

	t.Run("no comment", func(tt *testing.T) {
		const jsonc = `
		{
			"hello": "World"
		}
		`
		testJSON(tt, jsonc, &exp, &testData{})
	})

	t.Run("line comment", func(tt *testing.T) {
		const jsonc = `
		{
			// Test comment.
			"hello": "World"
		}
		`
		testJSON(tt, jsonc, &exp, &testData{})
	})

	t.Run("in-line comment", func(tt *testing.T) {
		const jsonc = `
		{
			"hello": "World" // Test comment.
		}
		`
		testJSON(tt, jsonc, &exp, &testData{})
	})

	t.Run("multiline comment at start", func(tt *testing.T) {
		const jsonc = `
		{
			/*
				Multiline
				comment.
			*/
			"hello": "World"
		}
		`
		testJSON(tt, jsonc, &exp, &testData{})
	})

	t.Run("multiline comment at end", func(tt *testing.T) {
		const jsonc = `
		{
			"hello": "World" /*	Multiline
				in-line
				comment.*/
		}
		`
		testJSON(tt, jsonc, &exp, &testData{})
	})

	t.Run("in-line comment at string", func(tt *testing.T) {
		const jsonc = `
		{
			"hello": "// world"
		}
		`

		var expData = testData{
			"// world",
		}

		testJSON(tt, jsonc, &expData, &testData{})
	})

	t.Run("multiline comment at string", func(tt *testing.T) {
		const jsonc = `
		{
			"hello": "/* world */"
		}
		`

		var expData = testData{
			"/* world */",
		}

		testJSON(tt, jsonc, &expData, &testData{})
	})
}

func TestComplexExample(t *testing.T) {
	const jsonc = `
	{ // ...
		/* ... */ "country_codes" /* ... */ : /* ... */ [ // ...
			{ // ...
				"country_code": "VN", // ...
				"country_name": "Vietnam", // ...
				"dialling_code": "+84" // ...
			}, // ...
			{ // ...
				"country_code": /* ... */ "JP", // ...
				"country_name": "Japan" /* ... */,
				"dialling_code": "+81" // ...
			}, /* ...
		*/	{
				"country_code": "TH" // ...
			,	"country_name": "Thailand" // ...
			// ...
			, /* ... */"dialling_code": "+66" // ...
			} // ...
		] // ...
	}// ...
	// ...
	`

	type countryCode struct {
		CountryCode  string `json:"country_code"`
		CountryName  string `json:"country_name"`
		DiallingCode string `json:"dialling_code"`
	}
	type countryCodesData struct {
		CountryCodes []countryCode `json:"country_codes"`
	}

	var exp = countryCodesData{
		CountryCodes: []countryCode{{
			CountryCode:  "VN",
			CountryName:  "Vietnam",
			DiallingCode: "+84",
		}, {
			CountryCode:  "JP",
			CountryName:  "Japan",
			DiallingCode: "+81",
		}, {
			CountryCode:  "TH",
			CountryName:  "Thailand",
			DiallingCode: "+66",
		}},
	}

	testJSON(t, jsonc, &exp, &countryCodesData{})
}

func testJSON(t *testing.T, in string, exp interface{}, got interface{}) {
	var r = jsoncjson.NewReader(strings.NewReader(in))

	var buff bytes.Buffer
	var _, err = buff.ReadFrom(r)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	t.Log(buff.String())

	err = json.NewDecoder(&buff).Decode(got)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(exp, got) {
		t.Fatalf("exp %+x, got: %+x", exp, got)
	}
}

func TestLargeJSON(t *testing.T) {
	var data = make(map[string]int, bytes.MinRead)
	for i := 0; i < bytes.MinRead; i++ {
		data[strconv.Itoa(i)] = i
	}

	var expBuf, err = json.Marshal(data)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var r = jsoncjson.NewReader(bytes.NewReader(expBuf))

	var gotBuf []byte
	gotBuf, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !bytes.Equal(expBuf, gotBuf) {
		t.Fatalf("exp len %d, got len: %d", len(expBuf), len(gotBuf))
	}
}

type errorReader struct {
	err error
}

func (r errorReader) Read([]byte) (n int, err error) {
	return 0, r.err
}

func TestThrownError(t *testing.T) {
	var expErr = errors.New("test error")

	var r = jsoncjson.NewReader(errorReader{
		err: expErr,
	})

	var _, err = r.Read(make([]byte, bytes.MinRead))

	if !errors.Is(err, expErr) {
		t.Fatalf("exp %s, got: %s", expErr, err)
	}
}
