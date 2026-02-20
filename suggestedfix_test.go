package logmsglint_test

import (
	"github.com/timur-developer/logmsglint"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestSuggestedFixes(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), logmsglint.Analyzer, "fix")
}
