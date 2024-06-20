//go:build ignore

package main

import (
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

// analyzer represents the checkcomments analyzer.
var analyzer = &analysis.Analyzer{
	Name: "getwd",
	Doc:  "getwd",
	URL:  "TODO",
	Run:  run,
}

// main runs the analyzer.
func main() {
	singlechecker.Main(analyzer)
}

// run analyses TODO comments.
func run(pass *analysis.Pass) (any, error) {
	wd, err := os.Getwd()
	pass.Reportf(0, "vet: wd = %q, err = %v", wd, err)
	return nil, nil
}
