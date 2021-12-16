package bson

import (
	"bufio"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AlekSi/go-bug/50214/types"
)

func parseDump(s string) ([]byte, error) {
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
			return nil, err
		}
		res = append(res, b...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func mustParseDumpFile(path ...string) []byte {
	b, err := os.ReadFile(filepath.Join(path...))
	if err != nil {
		panic(err)
	}

	b, err = parseDump(string(b))
	if err != nil {
		panic(err)
	}

	return b
}

var arrayTestCases = []testCase{{
	name: "array_all",
	v: &Array{
		types.Array{},
		types.MustMakeDocument(),
	},
	b: mustParseDumpFile("testdata", "array_all.hex"),
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
