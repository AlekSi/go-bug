package bson

import (
	"testing"
)

var cstringTestCases = []testCase{{
	b: []byte{0x66, 0x6f, 0x6f, 0x00},
}, {
	b: []byte{0x00},
}}

func FuzzCStringBinary(f *testing.F) {
	fuzzBinary(f, cstringTestCases)
}
