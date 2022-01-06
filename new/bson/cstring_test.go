package bson

import (
	"testing"

	"github.com/AlekSi/pointer"
)

var cstringTestCases = []testCase{{
	name: "foo",
	v:    pointer.To(CString("foo")),
	b:    []byte{0x66, 0x6f, 0x6f, 0x00},
	j:    `{"$c":"foo"}`,
}, {
	name: "empty",
	v:    pointer.To(CString("")),
	b:    []byte{0x00},
	j:    `{"$c":""}`,
}}

func TestCString(t *testing.T) {
	t.Parallel()

	t.Run("Binary", func(t *testing.T) {
		t.Parallel()
		testBinary(t, cstringTestCases, func() bsontype { return new(CString) })
	})

	t.Run("JSON", func(t *testing.T) {
		t.Parallel()
		testJSON(t, cstringTestCases, func() bsontype { return new(CString) })
	})
}

func FuzzCStringBinary(f *testing.F) {
	fuzzBinary(f, cstringTestCases, func() bsontype { return new(CString) })
}

func FuzzCStringJSON(f *testing.F) {
	fuzzJSON(f, cstringTestCases, func() bsontype { return new(CString) })
}
