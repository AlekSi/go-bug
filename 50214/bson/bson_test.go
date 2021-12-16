package bson

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlekSi/go-bug/50214/types"
)

type testCase struct {
	name string
	v    bsontype
	b    []byte
}

func testBinary(t *testing.T, testCases []testCase, newFunc func() bsontype) {
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.NotEmpty(t, tc.name, "name should not be empty")
			require.NotEmpty(t, tc.b, "b should not be empty")

			t.Parallel()

			t.Run("ReadFrom", func(t *testing.T) {
				t.Parallel()

				v := newFunc()
				br := bytes.NewReader(tc.b)
				bufr := bufio.NewReader(br)
				err := v.ReadFrom(bufr)
				assert.NoError(t, err)
				assert.Equal(t, tc.v, v, "expected: %s\nactual  : %s", tc.v, v)
				assert.Zero(t, br.Len(), "not all br bytes were consumed")
				assert.Zero(t, bufr.Buffered(), "not all bufr bytes were consumed")
			})

			t.Run("MarshalBinary", func(t *testing.T) {
				if tc.v == nil {
					t.Skip("v is nil")
				}

				t.Parallel()

				actualB, err := tc.v.MarshalBinary()
				require.NoError(t, err)
				if !assert.Equal(t, tc.b, actualB, "actual:\n%s", hex.Dump(actualB)) {
					// unmarshal again to compare BSON values
					v := newFunc()
					br := bytes.NewReader(actualB)
					bufr := bufio.NewReader(br)
					err := v.ReadFrom(bufr)
					assert.NoError(t, err)
					if assert.Equal(t, tc.v, v, "expected: %s\nactual  : %s", tc.v, v) {
						t.Log("BSON values are equal after unmarshalling")
					}
					assert.Zero(t, br.Len(), "not all br bytes were consumed")
					assert.Zero(t, bufr.Buffered(), "not all bufr bytes were consumed")
				}
			})

			t.Run("WriteTo", func(t *testing.T) {
				if tc.v == nil {
					t.Skip("v is nil")
				}

				t.Parallel()

				var buf bytes.Buffer
				bufw := bufio.NewWriter(&buf)
				err := tc.v.WriteTo(bufw)
				require.NoError(t, err)
				err = bufw.Flush()
				require.NoError(t, err)
				assert.Equal(t, tc.b, buf.Bytes(), "actual:\n%s", hex.Dump(buf.Bytes()))
			})
		})
	}
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

var bsonTestCases = []testCase{{
	name: "bson",
	v:    mustConvertDocument(types.MustMakeDocument()),
	b:    []byte{0x05, 0x00, 0x00, 0x00, 0x00},
}}

func TestBSON(t *testing.T) {
	t.Parallel()

	t.Run("Binary", func(t *testing.T) {
		t.Parallel()
		testBinary(t, bsonTestCases, func() bsontype { return new(Document) })
	})
}

func FuzzBSONBinary(f *testing.F) {
	fuzzBinary(f, bsonTestCases, func() bsontype { return new(Document) })
}
