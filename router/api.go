package router

import (
	. "github.com/projectriri/bot-gateway/model"
	"github.com/projectriri/bot-gateway/utils"
)

var producerBuffer = make(Buffer)
var producerChannelPool = make(map[string]*ProducerChannel)
var consumerChannelPool = make(map[string]*ConsumerChannel)

func RegisterProducerChannel(uuid string, requireAck bool) *ProducerChannel {
	uuid = utils.ValidateOrGenerateUUID(uuid)

	if pc, ok := producerChannelPool[uuid]; ok {
		return pc
	}

	var ackBuff *Buffer
	if requireAck {
		buff := make(Buffer)
		ackBuff = &buff
	}

	pc := &ProducerChannel{
		Channel: Channel{
			UUID:   uuid,
			Buffer: &producerBuffer,
		},
		AcknowlegeBuffer: ackBuff,
	}

	producerChannelPool[uuid] = pc
	return pc

}

func RegisterConsumerChannel(uuid string, accept []RoutingRule) *ConsumerChannel {
	uuid = utils.ValidateOrGenerateUUID(uuid)

	if cc, ok := consumerChannelPool[uuid]; ok {
		return cc
	}

	buff := make(Buffer)

	cc := &ConsumerChannel{
		Channel: Channel{
			UUID:   uuid,
			Buffer: &buff,
		},
		Accept: accept,
	}

	consumerChannelPool[uuid] = cc
	return cc
}

func Start() {

}
