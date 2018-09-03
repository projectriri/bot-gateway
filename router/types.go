package router

import "time"

type Packet struct {
	Head Head
	Body interface{}
}

type Head struct {
	UUID                  string
	From                  string
	To                    string
	ReplyToUUID           string
	AcknowlegeChannelUUID string
	Format                Format
}

type Format struct {
	API      string
	Version  string
	Protocol string
}

type Buffer chan Packet

type Channel struct {
	UUID       string
	Buffer     *Buffer
	ExpireTime time.Time
}

type ProducerChannel struct {
	Channel
	AcknowlegeBuffer *Buffer
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
