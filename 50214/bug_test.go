package bug

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

type Object struct{}

func (o *Object) ReadFrom(r *bufio.Reader) error {
	_, err := r.ReadByte()
	return err
}

func (o Object) WriteTo(w *bufio.Writer) error {
	_, err := w.Write([]byte{0x42})
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

		var v *Object
		var expectedB []byte

		// test ReadFrom
		{
			v = new(Object)
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
				t.Fatal("MarshalBinary results differ")
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
				t.Fatal("WriteTo results differ")
			}
		}
	})
}

func FuzzBinary(f *testing.F) {
	fuzzBinary(f, []testCase{{
		o: new(Object),
		b: []byte{0x42},
	}})
}
