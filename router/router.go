package router

import . "github.com/projectriri/bot-gateway/types"

var config RouterConfig
var producerBuffer Buffer
var producerChannelPool = make(map[string]*ProducerChannel)
var consumerChannelPool = make(map[string]*ConsumerChannel)
var converters []Converter

func Start(cfg RouterConfig, cvts []Converter) {
	config = cfg
	converters = cvts
	producerBuffer = make(Buffer, config.BufferSize)
	go garbageCollection()
	route()
}
