package router

import (
	. "github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/helpinfo"
	"regexp"
)

type Channel struct {
	UUID   string
	Buffer Buffer
}

type ProducerChannel struct {
	Channel
	AcknowledgeBuffer Buffer
}

type ConsumerChannel struct {
	Channel
	Accept   []RoutingRule
	HelpInfo *helpinfo.HelpInfo
}

type RoutingRule struct {
	From       string         `json:"from"`
	FromRegexp *regexp.Regexp `json:"-"`
	To         string         `json:"to"`
	ToRegexp   *regexp.Regexp `json:"-"`
	Formats    []Format       `json:"formats"`
}
