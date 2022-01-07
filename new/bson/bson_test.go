package bson

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

type testCase struct {
	name string
	v    bsontype
	b    []byte
}

func fuzzBinary(f *testing.F, testCases []testCase, newFunc func() bsontype) {
	for _, tc := range testCases {
		f.Add(tc.b)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		t.Parallel()

		var v bsontype
		var expectedB []byte

		// test ReadFrom
		{
			v = newFunc()
			br := bytes.NewReader(b)
			bufr := bufio.NewReader(br)
			if err := v.ReadFrom(bufr); err != nil {
				t.Skip(err)
			}

			// remove random tail
			expectedB = b[:len(b)-bufr.Buffered()-br.Len()]
		}

		// test MarshalBinary
		{
			actualB, err := v.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(expectedB, actualB) {
				t.Errorf("expected %v, got %v", expectedB, actualB)
			}
		}

		// test WriteTo
		{
			var bw bytes.Buffer
			bufw := bufio.NewWriter(&bw)
			err := v.WriteTo(bufw)
			if err != nil {
				t.Fatal(err)
			}
			err = bufw.Flush()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(expectedB, bw.Bytes()) {
				t.Errorf("expected %v, got %v", expectedB, bw.Bytes())
			}
		}
	})
}
