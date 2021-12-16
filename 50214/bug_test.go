package bug

import (
	"testing"
)

func FuzzParallel(f *testing.F) {
	f.Add(42)

	f.Fuzz(func(t *testing.T, i int) {
		t.Parallel()
	})
}
