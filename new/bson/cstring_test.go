package bson

import (
	"testing"
)

func cstringP(s string) *CString {
	c := CString(s)
	return &c
}

var cstringTestCases = []testCase{{
	b: []byte{0x66, 0x6f, 0x6f, 0x00},
}, {
	b: []byte{0x00},
}}

func FuzzCStringBinary(f *testing.F) {
	fuzzBinary(f, cstringTestCases, func() bsontype { return new(CString) })
}
