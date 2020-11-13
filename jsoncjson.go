// Package jsoncjson provides JSONC (JSON with comment) reader that
// removes all comments from the input.
//
// For example following input:
//
//	{ /* Comment. */
//		"Hello": "World" // Comment.
//	}
//
// Produces:
//
//	{
//		"Hello": "World"
//	}
//
package jsoncjson

import (
	"bytes"
	"errors"
	"io"
)

type jsoncTranslator struct {
	r io.Reader

	data   []byte
	end    bool
	cursor int

	lastByte  byte
	lastToken token
}

// NewReader creates new reader that removes all comments from a jsonc
// input and returns pure json.
func NewReader(r io.Reader) io.Reader {
	return &jsoncTranslator{
		r: r,

		data:   make([]byte, bytes.MinRead),
		cursor: bytes.MinRead,

		lastToken: tokenOther,
	}
}

// Read implements io.Reader interface.
func (t *jsoncTranslator) Read(jsonOut []byte) (n int, err error) {
	for n = range jsonOut {
		jsonOut[n], err = t.next()
		if err != nil {
			return n, err
		}
	}

	return n + 1, nil
}

type token int8

const (
	tokenString token = iota
	tokenSingleComment
	tokenMultiComment
	tokenUnknownComment
	tokenOther
	tokenEscaping
)

func (t *jsoncTranslator) handleToken(curByte byte) (skip bool) {
	switch t.lastToken {
	case tokenString:
		switch curByte {
		case '"':
			t.lastToken = tokenOther
		case '\\':
			t.lastToken = tokenEscaping
		}

		return false
	case tokenEscaping:
		t.lastToken = tokenString

		return false
	case tokenSingleComment:
		if curByte == '\n' {
			t.lastToken = tokenOther
		}

		return true
	case tokenMultiComment:
		if curByte == '/' && t.lastByte == '*' {
			t.lastToken = tokenOther
		}

		return true
	case tokenUnknownComment:
		switch curByte {
		case '/':
			t.lastToken = tokenSingleComment
		case '*':
			t.lastToken = tokenMultiComment
		}
		return true
	}

	switch curByte {
	case '"':
		t.lastToken = tokenString
	case '/':
		t.lastToken = tokenUnknownComment
		return true
	}

	return false
}

func (t *jsoncTranslator) refreshBuffer() (err error) {
	var n int
	n, err = t.r.Read(t.data)
	if err != nil {
		switch {
		case errors.Is(err, io.EOF):
			t.end = true
		default:
			return err
		}
	}

	t.data = t.data[:n]

	return nil
}

func (t *jsoncTranslator) next() (b byte, err error) {
	b, err = t.nextRawByte()
	if err != nil {
		return b, err
	}

	var skip = t.handleToken(b)
	if skip {
		t.lastByte = b
		return t.next()
	}

	t.lastByte = b
	return b, nil
}

func (t *jsoncTranslator) nextRawByte() (b byte, err error) {
	if t.cursor >= len(t.data) {
		if t.end {
			return 0, io.EOF
		}

		err = t.refreshBuffer()
		if err != nil {
			return 0, err
		}
		t.cursor = 0

		if len(t.data) == 0 {
			return 0, io.EOF
		}
	}

	b = t.data[t.cursor]
	t.cursor++

	return b, nil
}
