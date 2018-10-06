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
	Server      *Server
	channelPool map[string]*Channel
	expChan     chan bool
}

type Server struct {
	channelPool     map[string]*Channel
	gcInterval      time.Duration
	channelLifeTime time.Duration
}

func (s *Server) init(gci time.Duration, clt time.Duration) {
	s.channelPool = make(map[string]*Channel)
	s.gcInterval = gci
	s.channelLifeTime = clt
}

func (b *Broker) init(s *Server) {
	b.Server = s
	b.channelPool = make(map[string]*Channel)
	b.expChan = make(chan bool)
	go func() {
		for {
			select {
			case <-b.expChan:
				return
			case <-time.After(b.Server.channelLifeTime):
				for _, ch := range b.channelPool {
					b.Server.renewChannel(ch)
				}
			}
		}
	}()
}

func (b *Broker) close() {
	b.expChan <- false
	close(b.expChan)
}

func (b *Broker) InitChannel(args *ChannelInitRequest, reply *ChannelInitResponse) (err error) {
	uuid := utils.ValidateOrGenerateUUID(args.UUID)
	if !args.Producer && !args.Consumer {
		*reply = ChannelInitResponse{
			UUID: uuid,
			Code: 10042,
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
	b.Server.renewChannel(ch)
	b.channelPool[uuid] = ch
	b.Server.channelPool[uuid] = ch
	*reply = ChannelInitResponse{
		UUID: uuid,
		Code: 10001,
	}
	return
}

func (b *Broker) Send(args *ChannelProduceRequest, reply *ChannelProduceResponse) (err error) {
	log.Debugf("[jsonrpc-server-any] preparing to send packet")
	ch, ok := b.Server.channelPool[args.UUID]
	if !ok {
		*reply = ChannelProduceResponse{
			Code: 10044,
		}
		return
	}
	log.Debugf("[jsonrpc-server-any] sending packet for channel %v", ch.UUID)
	if ch.P == nil {
		*reply = ChannelProduceResponse{
			Code: 10048,
		}
		log.Warnf("[jsonrpc-server-any] sending packet for channel %v 418", ch.UUID)
		return
	}
	ch.P.Produce(args.Packet)
	*reply = ChannelProduceResponse{
		Code: 10002,
	}
	log.Debugf("[jsonrpc-server-any] sending packet for channel %v 202", ch.UUID)
	return
}

func (b *Broker) GetUpdates(args *ChannelConsumeRequest, reply *ChannelConsumeResponse) (err error) {
	log.Debugf("[jsonrpc-server-any] preparing to get updates")
	ch, ok := b.Server.channelPool[args.UUID]
	if !ok {
		*reply = ChannelConsumeResponse{
			Code: 10044,
		}
		return
	}
	log.Infof("[jsonrpc-server-any] preparing updates for channel %v", ch.UUID)
	if args.Limit <= 0 || args.Limit > 100 {
		args.Limit = 100
	}
	args.Timeout, err = time.ParseDuration(args.TimeoutStr)
	if err != nil {
		*reply = ChannelConsumeResponse{
			Code: 10040,
		}
		return
	}
	t := time.NewTimer(args.Timeout)
	for {
		select {
		case <-t.C:
			*reply = ChannelConsumeResponse{
				Code: 10004,
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
				Code:    10000,
				Packets: packets,
			}
			return
		}
	}
	return
}
