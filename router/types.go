package router

import (
	. "github.com/projectriri/bot-gateway/types"
	log "github.com/sirupsen/logrus"
	"time"
)

type Channel struct {
	UUID       string
	Buffer     Buffer
	ExpireTime time.Time
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
	From    string
	To      string
	Formats []Format
}

func getExpireTime() time.Time {
	return time.Now().Local().Add(config.ChannelLifeTime)
}

func (ch *Channel) renew() {
	ch.ExpireTime = getExpireTime()
	log.Debugf("[core] renewed channel %v to %v", ch.UUID, ch.ExpireTime)
}

func (pc ProducerChannel) Produce(packet Packet) {
	for {
		pc.renew()
		select {
		case pc.Buffer <- packet:
			log.Debugf("[core] pushed packet %v into buffer %v", packet.Head.UUID, pc.Buffer)
			return
		case <-time.After(config.ChannelLifeTime / 2):
		}
	}
}

func (cc ConsumerChannel) Consume() Packet {
	for {
		cc.renew()
		select {
		case packet := <-cc.Buffer:
			log.Debugf("[core] took packet %v from buffer %v", packet.Head.UUID, cc.Buffer)
			return packet
		case <-time.After(config.ChannelLifeTime / 2):
		}
	}
}
