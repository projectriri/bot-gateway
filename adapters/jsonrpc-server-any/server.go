package main

import (
	. "github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/jsonrpc-any"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/utils"
	log "github.com/sirupsen/logrus"
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
		*reply = ChannelInitResponse{
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
	*reply = ChannelInitResponse{
		UUID: uuid,
		Code: 201,
	}
	return
}

func (b *Broker) Send(args *ChannelProduceRequest, reply *ChannelProduceResponse) (err error) {
	ch, ok := b.channelPool[args.UUID]
	if !ok {
		*reply = ChannelProduceResponse{
			Code: 404,
		}
		return
	}
	b.renewChannel(ch)
	if ch.P == nil {
		*reply = ChannelProduceResponse{
			Code: 418,
		}
		return
	}
	ch.P.Produce(args.Packet)
	*reply = ChannelProduceResponse{
		Code: 200,
	}
	return
}

func (b *Broker) GetUpdates(args *ChannelConsumeRequest, reply *ChannelConsumeResponse) (err error) {
	log.Debugf("[jsonrpc-server-any] preparing to get updates")
	ch, ok := b.channelPool[args.UUID]
	if !ok {
		*reply = ChannelConsumeResponse{
			Code: 404,
		}
		return
	}
	b.renewChannel(ch)
	log.Infof("[jsonrpc-server-any] preparing updates for channel %v", ch.UUID)
	if args.Limit <= 0 || args.Limit > 100 {
		args.Limit = 100
	}
	t := time.NewTimer(args.Timeout)
	for {
		select {
		case <-t.C:
			*reply = ChannelConsumeResponse{
				Code: 204,
			}
			log.Infof("[jsonrpc-server-any] timeout")
			return
		case x := <-ch.C.Buffer:
			var packets []types.Packet
			packets = append(packets, x)
		More:
			for {
				if len(packets) >= args.Limit {
					break
				}
				select {
				case x = <-ch.C.Buffer:
					packets = append(packets, x)
				default:
					break More
				}
			}
			log.Debugf("[jsonrpc-server-any] %v", packets)
			*reply = ChannelConsumeResponse{
				Code:    200,
				Packets: packets,
			}
			return
		}
	}
	return
}
