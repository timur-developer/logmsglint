package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/timur-developer/logmsglint"
	"golang.org/x/tools/go/analysis"
)

type plugin struct{}

func init() {
	register.Plugin("logmsglint", New)
}

func New(settings any) (register.LinterPlugin, error) {
	return &plugin{}, nil
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{logmsglint.Analyzer}, nil
}

func (*plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
