package bson

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name string
	v    bsontype
	b    []byte
	bErr string // unwrapped

	j      string
	canonJ string // canonical form without extra object fields, zero values, etc.
	jErr   string // unwrapped
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

func testJSON(t *testing.T, testCases []testCase, newFunc func() bsontype) {
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.NotEmpty(t, tc.name, "name should not be empty")
			if tc.j == "" {
				t.Skip("j is empty")
			}

			t.Parallel()

			var dst bytes.Buffer
			require.NoError(t, json.Compact(&dst, []byte(tc.j)))
			require.Equal(t, tc.j, dst.String(), "j should be compacted")
			if tc.canonJ != "" {
				dst.Reset()
				require.NoError(t, json.Compact(&dst, []byte(tc.canonJ)))
				require.Equal(t, tc.canonJ, dst.String(), "canonJ should be compacted")
			}

			t.Run("UnmarshalJSON", func(t *testing.T) {
				t.Parallel()

				v := newFunc()
				err := v.UnmarshalJSON([]byte(tc.j))
				if tc.jErr == "" {
					require.NoError(t, err)
					assert.Equal(t, tc.v, v, "expected: %s\nactual  : %s", tc.v, v)
					return
				}

				require.Error(t, err)
				for {
					e := errors.Unwrap(err)
					if e == nil {
						break
					}
					err = e
				}
				require.Equal(t, tc.jErr, err.Error())
			})

			t.Run("MarshalJSON", func(t *testing.T) {
				t.Parallel()

				actualJ, err := tc.v.MarshalJSON()
				require.NoError(t, err)
				expectedJ := tc.j
				if tc.canonJ != "" {
					expectedJ = tc.canonJ
				}
				assert.Equal(t, expectedJ, string(actualJ))
			})
		})
	}
}

func fuzzJSON(f *testing.F, testCases []testCase, newFunc func() bsontype) {
	for _, tc := range testCases {
		f.Add(tc.j)
		if tc.canonJ != "" {
			f.Add(tc.canonJ)
		}
	}

	f.Fuzz(func(t *testing.T, j string) {
		t.Parallel()

		// raw "null" should never reach UnmarshalJSON due to the way encoding/json works
		if j == "null" {
			t.Skip(j)
		}

		// j may not be a canonical form.
		// We can't compare it with MarshalJSON() result directly.
		// Instead, we compare second results.

		v := newFunc()
		if err := v.UnmarshalJSON([]byte(j)); err != nil {
			t.Skip(err)
		}

		// test MarshalJSON
		{
			b, err := v.MarshalJSON()
			require.NoError(t, err)
			j = string(b)
		}

		// test UnmarshalJSON
		{
			actualV := newFunc()
			err := actualV.UnmarshalJSON([]byte(j))
			require.NoError(t, err)
			assert.Equal(t, v, actualV, "expected: %s\nactual  : %s", v, actualV)
		}
	})
}
