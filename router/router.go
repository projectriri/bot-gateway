package router

var config RouterConfig
var producerBuffer Buffer
var producerChannelPool = make(map[string]*ProducerChannel)
var consumerChannelPool = make(map[string]*ConsumerChannel)

func Start(cfg RouterConfig) {
	config = cfg
	producerBuffer = make(Buffer, config.BufferSize)
	go garbageCollection()
	route()
}
