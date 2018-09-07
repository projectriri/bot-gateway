package router

import . "github.com/projectriri/bot-gateway/types"

var config RouterConfig
var producerBuffer Buffer
var producerChannelPool = make(map[string]*ProducerChannel)
var consumerChannelPool = make(map[string]*ConsumerChannel)
var converters []Converter

func Init(cfg RouterConfig) {
	config = cfg
	producerBuffer = make(Buffer, config.BufferSize)
}

func Start(cvts []Converter) {
	converters = cvts
	route()
}
