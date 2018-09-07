package main

import (
	. "github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/jsonrpc-any"
	log "github.com/sirupsen/logrus"
	"time"
)

func (b *Broker) getExpireTime() time.Time {
	return time.Now().Local().Add(b.channelLifeTime)
}

func (b *Broker) renewChannel(ch *Channel) {
	ch.ExpireTime = b.getExpireTime()
	log.Debugf("[jsonrpc-server-any] renewed channel %v to %v", ch.UUID, ch.ExpireTime)
}

func (b *Broker) collectExpiredChannel() {
	lb := len(b.channelPool)
	for k, v := range b.channelPool {
		if v.ExpireTime.Before(time.Now()) {
			log.Warnf("[jsonrpc-server-any] [GC]: channel %v expired at %v", v.UUID, v.ExpireTime)
			if v.P != nil {
				v.P.Close()
			}
			if v.C != nil {
				v.C.Close()
			}
			delete(b.channelPool, k)
		}
	}
	lc := len(b.channelPool)
	log.Debugf("[jsonrpc-server-any] [GC]: channel %v -> %v", lb, lc)
}

func (b *Broker) garbageCollection() {
	for {
		<-time.After(b.gcInterval)
		b.collectExpiredChannel()
	}
}
