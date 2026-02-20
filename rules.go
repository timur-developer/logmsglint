package logmsglint

import (
	"strings"
	"unicode"
)

type issue struct {
	message   string
	suggested string
}

var currentConfig = defaultConfig()

var builtInKeywords = []string{
	"password",
	"passwd",
	"api_key",
	"apikey",
	"secret",
	"token",
	"authorization",
	"bearer",
}

func checkMessage(m extractedMessage) []issue {
	msg := normalizeSpaces(m.text)
	if msg == "" {
		return nil
	}

	var out []issue

	if iss, ok := ruleLowercaseStart(msg, currentConfig); ok {
		out = append(out, iss)
	}

	if iss, ok := ruleEnglishOnly(msg); ok {
		out = append(out, iss)
	}

	if iss, ok := ruleAllowedChars(msg, currentConfig); ok {
		out = append(out, iss)
	}

	if iss, ok := ruleSensitive(msg, currentConfig); ok {
		out = append(out, iss)
	}

	return out
}

func ruleLowercaseStart(msg string, cfg Config) (issue, bool) {
	runes := []rune(strings.TrimSpace(msg))
	if len(runes) == 0 {
		return issue{}, false
	}

	r := runes[0]
	if unicode.IsLetter(r) && unicode.IsUpper(r) {
		if !cfg.EnableFixes {
			return issue{
				message: "log message must start with a lowercase letter",
			}, true
		}
		runes[0] = unicode.ToLower(r)
		return issue{
			message:   "log message must start with a lowercase letter",
			suggested: string(runes),
		}, true
	}

	return issue{}, false
}

func ruleEnglishOnly(msg string) (issue, bool) {
	for _, r := range msg {
		if !unicode.IsLetter(r) {
			continue
		}

		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			continue
		}

		return issue{message: "log message must contain only English letters"}, true
	}

	return issue{}, false
}

func ruleAllowedChars(msg string, cfg Config) (issue, bool) {
	allowedPunct := cfg.AllowedPunct

	for _, r := range msg {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == ' ' {
			continue
		}
		if strings.ContainsRune(allowedPunct, r) {
			continue
		}
		return issue{message: "log message must not contain special characters or emoji"}, true
	}
	return issue{}, false
}

func ruleSensitive(msg string, cfg Config) (issue, bool) {
	l := strings.ToLower(msg)

	for _, kw := range builtInKeywords {
		if strings.Contains(l, kw) {
			return issue{message: "log message may contain sensitive data"}, true
		}
	}
	for _, kw := range cfg.ExtraKeywords {
		if strings.Contains(l, strings.ToLower(kw)) {
			return issue{message: "log message may contain sensitive data"}, true
		}
	}
	for _, re := range cfg.SensitivePatterns {
		if re.MatchString(msg) {
			return issue{message: "log message may contain sensitive data"}, true
		}
	}

	return issue{}, false
}
