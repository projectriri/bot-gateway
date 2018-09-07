package main

import (
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/utils"
	"time"
)

type Broker struct {
	channelPool     map[string]*Channel
	gcInterval      time.Duration
	channelLifeTime time.Duration
}

func (b *Broker) init(gci time.Duration, clt time.Duration) {
	b.channelPool = make(map[string]*Channel)
	b.gcInterval = gci
	b.channelLifeTime = clt
}

func (b *Broker) InitChannel(args *ChannelInitRequest, reply *ChannelInitResponse) (err error) {
	uuid := utils.ValidateOrGenerateUUID(args.UUID)
	if !args.Producer && !args.Consumer {
		reply = &ChannelInitResponse{
			UUID: uuid,
			Code: 204,
		}
	}
	ch := &Channel{
		UUID: uuid,
	}
	if args.Producer {
		ch.P = router.RegisterProducerChannel(uuid, true)
	}
	if args.Consumer {
		ch.C = router.RegisterConsumerChannel(uuid, args.Accept)
	}
	b.renewChannel(ch)
	b.channelPool[uuid] = ch
	reply = &ChannelInitResponse{
		UUID: uuid,
		Code: 201,
	}
	return nil
}

func (b *Broker) Send(args *ChannelProduceRequest, reply *ChannelProduceResponse) (err error) {
	return nil
}

func (b *Broker) GetUpdates(args *ChannelConsumeRequest, reply *ChannelConsumeResponse) (err error) {
	return nil
}
