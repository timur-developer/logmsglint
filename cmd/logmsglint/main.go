package main

import (
	"github.com/timur-developer/logmsglint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(logmsglint.Analyzer)
}
