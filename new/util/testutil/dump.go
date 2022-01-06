package testutil

import (
	"os"
	"path/filepath"

	"github.com/AlekSi/go-bug/new/util/hex"
)

func mustParseDump(s string) []byte {
	b, err := hex.ParseDump(s)
	if err != nil {
		panic(err)
	}
	return b
}

func MustParseDumpFile(path ...string) []byte {
	b, err := os.ReadFile(filepath.Join(path...))
	if err != nil {
		panic(err)
	}
	return mustParseDump(string(b))
}
