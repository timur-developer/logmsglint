package logmsglint

import (
	"go/ast"
	"go/constant"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"strconv"
	"strings"
)

type extractedMessage struct {
	text    string
	dynamic bool
	literal *ast.BasicLit // sets when msgExpr is string literal
}

func extractMessage(pass *analysis.Pass, exp ast.Expr) extractedMessage {
	if lit, ok := exp.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		if s, err := strconv.Unquote(lit.Value); err == nil {
			return extractedMessage{text: s, literal: lit}
		}
	}

	if tv, ok := pass.TypesInfo.Types[exp]; ok && tv.Value != nil && tv.Value.Kind() == constant.String {
		return extractedMessage{text: constant.StringVal(tv.Value)}
	}

	text, dynamic := collectStaticText(pass, exp)
	return extractedMessage{text: text, dynamic: dynamic}
}

func collectStaticText(pass *analysis.Pass, e ast.Expr) (string, bool) {
	switch x := e.(type) {
	case *ast.BasicLit:
		if x.Kind == token.STRING {
			s, err := strconv.Unquote(x.Value)
			return s, err != nil
		}
		return "", true

	case *ast.BinaryExpr:
		if x.Op != token.ADD {
			return "", true
		}

		ls, ld := collectStaticText(pass, x.X)
		rs, rd := collectStaticText(pass, x.Y)
		return ls + rs, ld || rd

	case *ast.CallExpr:
		if isFmtSprintf(pass, x) && len(x.Args) > 0 {
			m := extractMessage(pass, x.Args[0])
			if m.text != "" {
				return m.text, true
			}
		}
		return "", true
	default:
		return "", true
	}
}

func isFmtSprintf(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	obj := pass.TypesInfo.Uses[sel.Sel]
	if obj == nil || obj.Pkg() == nil {
		return false
	}
	return obj.Pkg().Path() == "fmt" && obj.Name() == "Sprintf"
}

func normalizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
