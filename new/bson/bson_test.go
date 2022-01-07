package bson

import (
	"testing"
)

type testCase struct {
	b []byte
}

func fuzzBinary(f *testing.F, testCases []testCase) {
	for _, tc := range testCases {
		f.Add(tc.b)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		t.Parallel()
	})
}
