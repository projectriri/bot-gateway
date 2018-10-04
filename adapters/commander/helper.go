package main

import "unicode"

const (
	ESCAPE_CHAR            = '\\'
	RESPONSE_CMD     uint8 = 1 << 0
	RESPONSE_CMDSTR  uint8 = 1 << 1
	RESPONSE_ARGS    uint8 = 1 << 2
	RESPONSE_ARGSTXT uint8 = 1 << 3
	RESPONSE_ARGSSTR uint8 = 1 << 4
)

func (p *CommanderPlugin) isWhiteChar(r rune) bool {
	return unicode.IsSpace(r)
}

func (p *CommanderPlugin) isQuoteChar(r rune) bool {
	return unicode.Is(unicode.Quotation_Mark, r)
}

func (p *CommanderPlugin) checkPrefix(text string) (prefix string, containPrefix bool) {
	if len(p.config.CommandPrefix) == 0 {
		return "", true
	}
	for _, pfx := range p.config.CommandPrefix {
		if len(text) < len(pfx) {
			continue
		}
		if text[:len(pfx)] == pfx {
			return pfx, true
		}
	}
	return "", false
}

func (p *CommanderPlugin) checkAllowEmptyPrefix() bool {
	if len(p.config.CommandPrefix) == 0 {
		return true
	}
	for _, pfx := range p.config.CommandPrefix {
		if pfx == "" {
			return true
		}
	}
	return false
}
