package bug

import (
	"bufio"
	"testing"
)

type Object struct{}

func (o *Object) ReadFrom(r *bufio.Reader) error {
	_, err := r.ReadByte()
	return err
}

func (o Object) MarshalBinary() ([]byte, error) {
	return []byte{0x42}, nil
}

type testCase struct {
	o *Object
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

func FuzzBinary(f *testing.F) {
	fuzzBinary(f, []testCase{{
		o: new(Object),
		b: []byte{0x42},
	}})
}
