package router

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func collectExpiredChannel() {
	lb := len(producerChannelPool)
	for k, v := range producerChannelPool {
		if v.ExpireTime.Before(time.Now()) {
			delete(producerChannelPool, k)
		}
	}
	lc := len(producerChannelPool)
	log.Debugf("[GC]: Producer %v -> %v", lb, lc)
	lb = len(consumerChannelPool)
	for k, v := range consumerChannelPool {
		if v.ExpireTime.Before(time.Now()) {
			delete(consumerChannelPool, k)
		}
	}
	lc = len(consumerChannelPool)
	log.Debugf("[GC]: Consumer %v -> %v", lb, lc)
}

func garbageCollection() {
	for {
		t := time.NewTimer(config.GCInterval)
		<-t.C
		collectExpiredChannel()
	}
}
