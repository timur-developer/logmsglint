package logmsglint_test

import (
	"github.com/timur-developer/logmsglint"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), logmsglint.Analyzer, "lintcases")
}
