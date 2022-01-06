package bson

import (
	"testing"

	"github.com/AlekSi/go-bug/new/types"
)

var arrayTestCases = []testCase{{
	name: "array_all",
	v: &Array{
		types.Array{},
		types.MustMakeDocument(),
	},
	b: []byte{0x15, 0x00, 0x00, 0x00, 0x04, 0x30, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x03, 0x31, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00},
	j: "[[],{\"$k\":[]}]",
}}

func FuzzArrayBinary(f *testing.F) {
	fuzzBinary(f, arrayTestCases, func() bsontype { return new(Array) })
}
