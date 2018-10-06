package main

import (
	. "github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/jsonrpc-any"
	log "github.com/sirupsen/logrus"
	"time"
)

func (s *Server) getExpireTime() time.Time {
	return time.Now().Local().Add(s.channelLifeTime)
}

func (s *Server) renewChannel(ch *Channel) {
	ch.ExpireTime = s.getExpireTime()
	log.Debugf("[jsonrpc-server-any] renewed channel %v to %v", ch.UUID, ch.ExpireTime)
}

func (s *Server) collectExpiredChannel() {
	lb := len(s.channelPool)
	for k, v := range s.channelPool {
		if v.ExpireTime.Before(time.Now()) {
			log.Warnf("[jsonrpc-server-any] [GC]: channel %v expired at %v", v.UUID, v.ExpireTime)
			if v.P != nil {
				v.P.Close()
			}
			if v.C != nil {
				v.C.Close()
			}
			delete(s.channelPool, k)
		}
	}
	lc := len(s.channelPool)
	log.Debugf("[jsonrpc-server-any] [GC]: channel %v -> %v", lb, lc)
}

func (s *Server) garbageCollection() {
	for {
		<-time.After(s.gcInterval)
		s.collectExpiredChannel()
	}
}
