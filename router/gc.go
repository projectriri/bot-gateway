package router

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func collectExpiredChannel() {
	lb := len(producerChannelPool)
	for k, v := range producerChannelPool {
		if v.ExpireTime.Before(time.Now()) {
			log.Warnf("[GC]: Producer %v expired at %v", v.ExpireTime)
			delete(producerChannelPool, k)
		}
	}
	lc := len(producerChannelPool)
	log.Debugf("[GC]: Producer %v -> %v", lb, lc)
	lb = len(consumerChannelPool)
	for k, v := range consumerChannelPool {
		if v.ExpireTime.Before(time.Now()) {
			log.Warnf("[GC]: Consumer %v expired at %v", v.ExpireTime)
			delete(consumerChannelPool, k)
		}
	}
	lc = len(consumerChannelPool)
	log.Debugf("[GC]: Consumer %v -> %v", lb, lc)
}

func garbageCollection() {
	for {
		<-time.After(config.GCInterval)
		collectExpiredChannel()
	}
}
