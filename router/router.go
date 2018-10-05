package router

import (
	. "github.com/projectriri/bot-gateway/types"
)

var config RouterConfig
var producerBuffer Buffer
var producerChannelPool = make(map[string]*ProducerChannel)
var consumerChannelPool = make(map[string]*ConsumerChannel)
var converters []Converter
var inCnt = 0
var outCnt = 0

func Init(cfg RouterConfig) {
	config = cfg
	producerBuffer = make(Buffer, config.BufferSize)
}

func Start(cvts []Converter) {
	converters = cvts
	route()
}

func GetProducerChannelCount() int {
	return len(producerChannelPool)
}

func GetConsumerChannelCount() int {
	return len(consumerChannelPool)
}

func GetChannelCount() int {
	return GetProducerChannelCount() + GetConsumerChannelCount()
}

func GetCachedPacketCount() int {
	cnt := len(producerBuffer)
	for _, cc := range consumerChannelPool {
		cnt += len(cc.Buffer)
	}
	return cnt
}

func GetChannelCacheLimit() int {
	return int(config.BufferSize)
}

func GetIOCount() (int, int) {
	return inCnt, outCnt
}
