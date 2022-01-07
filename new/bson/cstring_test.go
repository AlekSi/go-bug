package bson

import (
	"testing"
)

func cstringP(s string) *CString {
	c := CString(s)
	return &c
}

var cstringTestCases = []testCase{{
	name: "foo",
	v:    cstringP("foo"),
	b:    []byte{0x66, 0x6f, 0x6f, 0x00},
}, {
	name: "empty",
	v:    cstringP(""),
	b:    []byte{0x00},
}}

func FuzzCStringBinary(f *testing.F) {
	fuzzBinary(f, cstringTestCases, func() bsontype { return new(CString) })
}
