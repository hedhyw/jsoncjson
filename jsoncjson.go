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

	lastByte byte
	curToken token
}

// NewReader creates new reader that removes all comment from the jsonc
// input and returns pure json.
func NewReader(r io.Reader) io.Reader {
	return &jsoncTranslator{
		r: r,

		data:   make([]byte, bytes.MinRead),
		cursor: bytes.MinRead,

		curToken: tokenOther,
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

	return n, nil
}

func (t *jsoncTranslator) handleToken(curByte byte) (skip bool) {
	switch t.curToken {
	case tokenString:
		if curByte == '"' {
			t.curToken = tokenOther
		}

		return false
	case tokenSingleComment:
		if curByte == '\n' {
			t.curToken = tokenOther
		}

		return true
	case tokenMultiComment:
		if curByte == '/' && t.lastByte == '*' {
			t.curToken = tokenOther
		}

		return true
	case tokenUnknownComment:
		switch curByte {
		case '/':
			t.curToken = tokenSingleComment
		case '*':
			t.curToken = tokenMultiComment
		}
		return true
	}

	switch curByte {
	case '"':
		t.curToken = tokenString
	case '/':
		t.curToken = tokenUnknownComment
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
