package router

import (
	. "github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/utils"
)

func RegisterProducerChannel(uuid string, acceptAck bool) *ProducerChannel {
	uuid = utils.ValidateOrGenerateUUID(uuid)

	if pc, ok := producerChannelPool[uuid]; ok {
		return pc
	}

	var ackBuff Buffer
	if acceptAck {
		buff := make(Buffer, config.BufferSize)
		ackBuff = buff
	}

	pc := &ProducerChannel{
		Channel: Channel{
			UUID:   uuid,
			Buffer: producerBuffer,
		},
		AcknowledgeBuffer: ackBuff,
	}

	producerChannelPool[uuid] = pc
	return pc

}

func RegisterConsumerChannel(uuid string, accept []RoutingRule) *ConsumerChannel {
	uuid = utils.ValidateOrGenerateUUID(uuid)

	if cc, ok := consumerChannelPool[uuid]; ok {
		return cc
	}

	buff := make(Buffer, config.BufferSize)

	cc := &ConsumerChannel{
		Channel: Channel{
			UUID:   uuid,
			Buffer: buff,
		},
		Accept: accept,
	}

	consumerChannelPool[uuid] = cc
	return cc
}

func (pc *ProducerChannel) Produce(packet Packet) {
	pc.Buffer <- packet
}

func (cc *ConsumerChannel) Consume() Packet {
	packet := <-cc.Buffer
	return packet
}

func (pc *ProducerChannel) Close() {
	delete(producerChannelPool, pc.UUID)
	pc = nil
}

func (cc *ConsumerChannel) Close() {
	delete(consumerChannelPool, cc.UUID)
	cc = nil
}
