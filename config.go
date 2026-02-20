package logmsglint

import (
	"fmt"
	"regexp"
)

type Config struct {
	EnableFixes       bool
	AllowedPunct      string
	ExtraKeywords     []string
	SensitivePatterns []*regexp.Regexp
}

type rawConfig struct {
	EnableFixes       bool     `yaml:"enable_fixes"`
	AllowedPunct      string   `yaml:"allowed_punct"`
	SensitiveKeywords []string `yaml:"sensitive_keywords"`
	SensitiveRegexps  []string `yaml:"sensitive_regexps"`
}

func defaultConfig() Config {
	return Config{
		EnableFixes:  true,
		AllowedPunct: "",
	}
}

func loadConfig(settings any) (Config, error) {
	cfg := defaultConfig()
	if settings == nil {
		return cfg, nil
	}

	m, ok := settings.(map[string]any)
	if !ok {
		return cfg, fmt.Errorf("logmsglint: unexpected settings type %T", settings)
	}

	var raw rawConfig

	if v, ok := m["enable_fixes"].(bool); ok {
		raw.EnableFixes = v
	}
	if v, ok := m["allowed_punct"].(string); ok {
		raw.AllowedPunct = v
	}
	if v, ok := m["sensitive_keywords"].([]any); ok {
		for _, x := range v {
			if s, ok := x.(string); ok {
				raw.SensitiveKeywords = append(raw.SensitiveKeywords, s)
			}
		}
	}
	if v, ok := m["sensitive_regexps"].([]any); ok {
		for _, x := range v {
			if s, ok := x.(string); ok {
				raw.SensitiveRegexps = append(raw.SensitiveRegexps, s)
			}
		}
	}

	cfg.EnableFixes = raw.EnableFixes
	cfg.AllowedPunct = raw.AllowedPunct
	cfg.ExtraKeywords = raw.SensitiveKeywords

	for _, pattern := range raw.SensitiveRegexps {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return cfg, fmt.Errorf("logmsglint: invalid regexp %q: %w", pattern, err)
		}
		cfg.SensitivePatterns = append(cfg.SensitivePatterns, re)
	}

	return cfg, nil
}

func LoadConfig(settings any) (Config, error) {
	return loadConfig(settings)
}

func SetConfig(cfg Config) {
	currentConfig = cfg
}
