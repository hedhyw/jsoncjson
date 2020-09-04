package jsoncjson

type token int8

const (
	tokenString token = iota
	tokenSingleComment
	tokenMultiComment
	tokenUnknownComment
	tokenOther
)
