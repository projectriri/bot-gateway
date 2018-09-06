package router

import (
	. "github.com/projectriri/bot-gateway/types"
	"time"
)

type Channel struct {
	UUID       string
	Buffer     *Buffer
	ExpireTime time.Time
}

type ProducerChannel struct {
	Channel
	AcknowledgeBuffer *Buffer
}

type ConsumerChannel struct {
	Channel
	Accept []RoutingRule
}

type RoutingRule struct {
	From    string
	To      string
	Formats []Format
}

func getExpireTime() time.Time {
	return time.Now().Local().Add(config.ChannelLifeTime)
}

func (ch Channel) renew() {
	ch.ExpireTime = getExpireTime()
}

func (pc ProducerChannel) Produce(packet Packet) {
	for {
		pc.renew()
		select {
		case *pc.Buffer <- packet:
			return
		case <-time.After(config.ChannelLifeTime / 2):
		}
	}
}

func (cc ConsumerChannel) Consume() Packet {
	for {
		cc.renew()
		select {
		case packet := <-*cc.Buffer:
			return packet
		case <-time.After(config.ChannelLifeTime / 2):
		}
	}
}
