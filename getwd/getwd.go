package main

import (
	"fmt"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var analyzer = &analysis.Analyzer{
	Name: "getwd",
	Doc:  "getwd",
	URL:  "TODO",
	Run: func(pass *analysis.Pass) (any, error) {
		wd, _ := os.Getwd()
		fmt.Printf("os.Getwd returned %s\n", wd)
		fmt.Printf("os.Getwd returned %q\n", wd)
		return nil, nil
	},
}

func main() {
	singlechecker.Main(analyzer)
}
