package main

import (
	"os"
	"testing"
)

func TestGetwd(t *testing.T) {
	wd, err := os.Getwd()
	t.Logf("test: wd = %s, err = %v", wd, err)
}
