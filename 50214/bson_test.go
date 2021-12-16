package bug

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			require.NoError(t, err)
			assert.Equal(t, expectedB, actualB, "MarshalBinary results differ")
		}

		// test WriteTo
		{
			var bw bytes.Buffer
			bufw := bufio.NewWriter(&bw)
			err := v.WriteTo(bufw)
			require.NoError(t, err)
			err = bufw.Flush()
			require.NoError(t, err)
			assert.Equal(t, expectedB, bw.Bytes(), "WriteTo results differ")
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
