package bson

import (
	"testing"

	"github.com/AlekSi/go-bug/50214/types"
	"github.com/AlekSi/go-bug/50214/util/testutil"
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

func TestArray(t *testing.T) {
	t.Parallel()

	t.Run("Binary", func(t *testing.T) {
		t.Parallel()
		testBinary(t, arrayTestCases, func() bsontype { return new(Array) })
	})

	t.Run("JSON", func(t *testing.T) {
		t.Parallel()
		testJSON(t, arrayTestCases, func() bsontype { return new(Array) })
	})
}

func FuzzArrayBinary(f *testing.F) {
	fuzzBinary(f, arrayTestCases, func() bsontype { return new(Array) })
}

func FuzzArrayJSON(f *testing.F) {
	fuzzJSON(f, arrayTestCases, func() bsontype { return new(Array) })
}

func BenchmarkArray(b *testing.B) {
	benchmark(b, arrayTestCases, func() bsontype { return new(Array) })
}
