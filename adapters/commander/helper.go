package main

import "unicode"

const ESCAPE_CHAR = '\\'

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
