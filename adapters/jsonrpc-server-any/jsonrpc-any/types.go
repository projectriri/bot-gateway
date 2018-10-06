package jsonrpc_any

import (
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"time"
)

type UpdatesChannel <-chan *types.Packet

type Channel struct {
	UUID       string
	ExpireTime time.Time
	P          *router.ProducerChannel
	C          *router.ConsumerChannel
}

type ChannelInitRequest struct {
	UUID     string               `json:"uuid"`
	Producer bool                 `json:"producer"`
	Consumer bool                 `json:"consumer"`
	Accept   []router.RoutingRule `json:"accept"`
}

type ChannelInitResponse struct {
	UUID string `json:"uuid"`
	Code int    `json:"code"`
}

type ChannelConsumeRequest struct {
	UUID       string        `json:"uuid"`
	Timeout    time.Duration `json:"-"`
	TimeoutStr string        `json:"timeout"`
	Limit      int           `json:"limit"`
}

type ChannelConsumeResponse struct {
	Packets []types.Packet `json:"packets"`
	Code    int            `json:"code"`
}

type ChannelProduceRequest struct {
	UUID   string       `json:"uuid"`
	Packet types.Packet `json:"packet"`
}

type ChannelProduceResponse struct {
	Code int `json:"code"`
}
