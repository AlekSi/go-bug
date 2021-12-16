package bson

import (
	"bufio"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/AlekSi/go-bug/50214/types"
)

func mustParseDump(s string) []byte {
	var res []byte

	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(s)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if line[len(line)-1] == '|' {
			// go dump
			line = strings.TrimSpace(line[8:60])
			line = strings.Join(strings.Split(line, " "), "")
		} else {
			// wireshark dump
			line = strings.TrimSpace(line[7:54])
			line = strings.Join(strings.Split(line, " "), "")
		}

		b, err := hex.DecodeString(line)
		if err != nil {
			panic(err)
		}
		res = append(res, b...)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return res
}

var arrayTestCases = []testCase{{
	name: "array_all",
	v: &Array{
		types.Array{},
		types.MustMakeDocument(),
	},
	b: mustParseDump(`
	00000000  15 00 00 00 04 30 00 05  00 00 00 00 03 31 00 05  |.....0.......1..|
	00000010  00 00 00 00 00                                    |.....|
		`),
	j: "[[],{\"$k\":[]}]",
}}

func TestArray(t *testing.T) {
	t.Parallel()

	t.Run("Binary", func(t *testing.T) {
		t.Parallel()
		testBinary(t, arrayTestCases, func() bsontype { return new(Array) })
	})
}

func FuzzArrayBinary(f *testing.F) {
	fuzzBinary(f, arrayTestCases, func() bsontype { return new(Array) })
}
