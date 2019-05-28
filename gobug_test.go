package gobug

import (
	"testing"
	"time"
)

type Params struct {
	L *time.Location
}

type Struct struct {
	Str1 string
	Str2 string
	Str3 string
	Str4 string
	S1   []string
	S2   []string
	M1   map[string]float64
	M2   map[string]uint64
	M3   map[string]bool
	U    uint
}

func f(t *testing.T, opt Params) []Struct {
	res := []Struct{}
	for {
		e := (*Struct)(nil)
		if e != nil {
			return res
		}
		res = append(res, *e)
	}
}

func TestGoBug(t *testing.T) {
	f(t, Params{L: time.UTC})
}
