package logmsglint

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
)

var Analyzer = &analysis.Analyzer{
	Name:     "logmsglint",
	Doc:      "checks slog/go.uber.zap log messages to basic style and safety rules",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

type callInfo struct {
	ok      bool
	msgExpr ast.Expr
}

// take only log functions from slog
var slogMsgIndex = map[string]int{
	"Debug": 0,
	"Info":  0,
	"Warn":  0,
	"Error": 0,

	"DebugContext": 1,
	"InfoContext":  1,
	"WarnContext":  1,
	"ErrorContext": 1,
}

// take only log functions from go.uber.zap
var zapMsgIndex = map[string]int{
	"Debug":  0,
	"Info":   0,
	"Warn":   0,
	"Error":  0,
	"DPanic": 0,
	"Panic":  0,
	"Fatal":  0,
}

func run(pass *analysis.Pass) (any, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}
	ins.Preorder(nodeFilter, func(node ast.Node) {
		call := node.(*ast.CallExpr)

		info := matchLoggerCall(pass, call)
		if !info.ok {
			return
		}

		if info.msgExpr == nil {
			return
		}

		extracted := extractMessage(pass, info.msgExpr)
		if extracted.text == "" {
			return
		}

		issues := checkMessage(extracted)
		if len(issues) == 0 {
			return
		}

		msgs := make([]string, 0, len(issues))
		suggested := ""

		for _, iss := range issues {
			msgs = append(msgs, iss.message)
			if suggested == "" && iss.suggested != "" {
				suggested = iss.suggested
			}
		}

		diag := analysis.Diagnostic{
			Pos:     info.msgExpr.Pos(),
			End:     info.msgExpr.End(),
			Message: strings.Join(msgs, "; "),
		}

		if currentConfig.EnableFixes && suggested != "" && extracted.literal != nil {
			if edit, ok := suggestedFixForStringLiteral(extracted.literal, suggested); ok {
				diag.SuggestedFixes = []analysis.SuggestedFix{{
					Message:   "apply fix",
					TextEdits: []analysis.TextEdit{edit},
				}}
			}
		}

		pass.Report(diag)
	})

	return nil, nil
}

func matchLoggerCall(pass *analysis.Pass, call *ast.CallExpr) callInfo {
	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		name := fun.Sel.Name

		if pkgPath := receiverPkgPath(pass, fun.X); pkgPath != "" {
			switch pkgPath {
			case "log/slog":
				idx, ok := slogMsgIndex[name]
				if !ok || idx >= len(call.Args) || !isString(pass, call.Args[idx]) {
					return callInfo{}
				}
				return callInfo{ok: true, msgExpr: call.Args[idx]}

			case "go.uber.org/zap":
				idx, ok := zapMsgIndex[name]
				if !ok || idx >= len(call.Args) || !isString(pass, call.Args[idx]) {
					return callInfo{}
				}
				return callInfo{ok: true, msgExpr: call.Args[idx]}
			}
		}

		if obj := pass.TypesInfo.Uses[fun.Sel]; obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == "log/slog" {
			idx, ok := slogMsgIndex[name]
			if !ok || idx >= len(call.Args) || !isString(pass, call.Args[idx]) {
				return callInfo{}
			}
			return callInfo{ok: true, msgExpr: call.Args[idx]}
		}

		return callInfo{}
	case *ast.Ident:
		obj := pass.TypesInfo.Uses[fun]
		if obj == nil || obj.Pkg() == nil || obj.Pkg().Path() != "log/slog" {
			return callInfo{}
		}

		idx, ok := slogMsgIndex[fun.Name]
		if !ok || idx >= len(call.Args) || !isString(pass, call.Args[idx]) {
			return callInfo{}
		}
		return callInfo{ok: true, msgExpr: call.Args[idx]}
	}

	return callInfo{}
}

func receiverPkgPath(pass *analysis.Pass, x ast.Expr) string {
	t := pass.TypesInfo.TypeOf(x)
	if t == nil {
		return ""
	}

	for {
		p, ok := t.(*types.Pointer)
		if !ok {
			break
		}
		t = p.Elem()
	}

	n, ok := t.(*types.Named)
	if !ok || n.Obj() == nil || n.Obj().Pkg() == nil {
		return ""
	}

	return n.Obj().Pkg().Path()
}

func isString(pass *analysis.Pass, e ast.Expr) bool {
	tv, ok := pass.TypesInfo.Types[e]
	if !ok || tv.Type == nil {
		return false
	}

	return types.Identical(tv.Type, types.Typ[types.String])
}
