package router

import (
	"time"
)

func collectExpiredChannel() {
	//lb := len(ChannelPool)
	for k, v := range producerChannelPool {
		if v.ExpireTime.Before(time.Now()) {
			delete(producerChannelPool, k)
		}
	}
	//lc := len(ChannelPool)
	//log.Debugf("[GC]: %v -> %v", lb, lc)
	//lb := len(ChannelPool)
	for k, v := range consumerChannelPool {
		if v.ExpireTime.Before(time.Now()) {
			delete(consumerChannelPool, k)
		}
	}
	//lc := len(ChannelPool)
	//log.Debugf("[GC]: %v -> %v", lb, lc)
}

func garbageCollection() {
	for {
		t := time.NewTimer(config.GCInterval)
		<-t.C
		collectExpiredChannel()
	}
}
