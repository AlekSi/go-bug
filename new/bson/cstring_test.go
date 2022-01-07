package bson

import (
	"testing"
)

func Fuzz1(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, b []byte) {
		t.Parallel()
	})
}
