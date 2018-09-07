package router

import (
	. "github.com/projectriri/bot-gateway/types"
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
	From    string   `json:"from"`
	To      string   `json:"to"`
	Formats []Format `json:"formats"`
}
