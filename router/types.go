package router

import (
	. "github.com/projectriri/bot-gateway/types"
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
	Accept []RoutingRule
}

type RoutingRule struct {
	From       string         `json:"from"`
	FromRegexp *regexp.Regexp `json:"-"`
	To         string         `json:"to"`
	ToRegexp   *regexp.Regexp `json:"-"`
	Formats    []Format       `json:"formats"`
}
