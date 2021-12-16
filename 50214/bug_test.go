package bug

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

type Document struct {
	m    map[string]any
	keys []string
}

func (doc *Document) ReadFrom(r *bufio.Reader) error {
	r.Read(make([]byte, 5))

	doc.m = map[string]any{}
	doc.keys = []string{}
	return nil
}

func (doc Document) WriteTo(w *bufio.Writer) error {
	_, err := w.Write([]byte{0x05, 0x00, 0x00, 0x00, 0x00})
	return err
}

func (doc Document) MarshalBinary() ([]byte, error) {
	return []byte{0x05, 0x00, 0x00, 0x00, 0x00}, nil
}

type testCase struct {
	name string
	v    *Document
	b    []byte
}

func fuzzBinary(f *testing.F, testCases []testCase) {
	for _, tc := range testCases {
		f.Add(tc.b)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		t.Parallel()

		var v *Document
		var expectedB []byte

		// test ReadFrom
		{
			v = new(Document)
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

func FuzzBSONBinary(f *testing.F) {
	bsonTestCases := []testCase{{
		name: "bson",
		v: &Document{
			m:    make(map[string]any),
			keys: []string{},
		},
		b: []byte{0x05, 0x00, 0x00, 0x00, 0x00},
	}}

	fuzzBinary(f, bsonTestCases)
}
