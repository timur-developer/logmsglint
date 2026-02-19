package logmsglint

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strconv"
	"strings"
)

func suggestedFixForStringLiteral(lit *ast.BasicLit, newValue string) (analysis.TextEdit, bool) {
	if lit == nil {
		return analysis.TextEdit{}, false
	}

	repl := strconv.Quote(newValue)

	// if it was raw-string we save quotes
	if strings.HasPrefix(lit.Value, "`") && !strings.ContainsRune(newValue, '`') {
		repl = "`" + newValue + "`"
	}

	return analysis.TextEdit{
		Pos:     lit.Pos(),
		End:     lit.End(),
		NewText: []byte(repl),
	}, true

}
