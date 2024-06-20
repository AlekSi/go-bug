//go:build ignore

package main

import (
	"fmt"
	"os"
)

func main() {
	wd, err := os.Getwd()
	fmt.Printf("main: wd = %q, err = %v\n", wd, err)
}
