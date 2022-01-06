package bson

import (
	"testing"

	"github.com/AlekSi/go-bug/new/types"
	"github.com/AlekSi/go-bug/new/util/testutil"
)

var arrayTestCases = []testCase{{
	name: "array_all",
	v: &Array{
		types.Array{},
		types.MustMakeDocument(),
	},
	b: testutil.MustParseDumpFile("testdata", "array_all.hex"),
	j: "[[],{\"$k\":[]}]",
}}

func FuzzArrayBinary(f *testing.F) {
	fuzzBinary(f, arrayTestCases, func() bsontype { return new(Array) })
}
