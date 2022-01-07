package new

import (
	"testing"
)

func Fuzz1(f *testing.F) {
	f.Add(42)
	f.Add(43)

	f.Fuzz(func(t *testing.T, i int) {
		t.Parallel()
	})
}
