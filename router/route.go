package router

import (
	. "github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/utils"
	log "github.com/sirupsen/logrus"
	"sync"
)

func route() {
	for {
		pkt := <-producerBuffer
		go processPacket(&pkt)
	}
}

func processPacket(pkt *Packet) {
	if !utils.ValidateUUID(pkt.Head.UUID) {
		log.Warnf("[router] pkt with invalid uuid, dropped: %+v BODY: %s", pkt.Head, string(pkt.Body))
		return
	}

	inCnt++
	log.Debugf("[router] pkt: %+v", pkt.Head)
	log.Debugf("[router] pkt: %s BODY: %s", pkt.Head.UUID, string(pkt.Body))

	processedMap := make(map[string][]Packet)
	processingMap := make(map[string][]chan []Packet)
	mux := sync.Mutex{}
	for _, consumerChannel := range consumerChannelPool {
		go func(cc *ConsumerChannel) {
			log.Debugf("[router] pkt: %v TRYING cc: %+v", pkt.Head.UUID, cc)
			var formats []Format
			for _, ac := range cc.Accept {
				f := ac.FromRegexp.MatchString(pkt.Head.From)
				t := ac.ToRegexp.MatchString(pkt.Head.To)
				if f && t {
					formats = ac.Formats
					break
				}
			}

			if formats == nil || len(formats) == 0 {
				return
			}

			log.Debugf("[router] pkt: %v ACCEPTED BY cc: %+v", pkt.Head.UUID, cc)

			for _, format := range formats {
				f := format.String()

				mux.Lock()
				if convertedPacket, ok := processedMap[f]; ok {
					if convertedPacket == nil {
						mux.Unlock()
						continue
					}
					pushPacket(cc, convertedPacket)
					mux.Unlock()
					return
				}

				if pkt.Head.Format.String() == f {
					processedMap[f] = []Packet{*pkt}
					pushPacket(cc, processedMap[f])
					mux.Unlock()
					return
				}

				if converting, ok := processingMap[f]; ok && converting != nil {
					ch := make(chan []Packet)
					converting = append(converting, ch)
					mux.Unlock()
					result := <-ch
					if result != nil {
						pushPacket(cc, result)
						return
					}
					continue
				}

				processingMap[f] = make([]chan []Packet, 0)
				mux.Unlock()
				for _, cvt := range converters {
					if cvt.IsConvertible(pkt.Head.Format, format) {
						ok, result := cvt.Convert(*pkt, format)
						if ok && result != nil {
							pushPacket(cc, result)
							mux.Lock()
							for _, ch := range processingMap[f] {
								ch <- result
								close(ch)
							}
							delete(processingMap, f)
							processedMap[f] = result
							mux.Unlock()
							return
						}
					}
				}
				mux.Lock()
				for _, ch := range processingMap[f] {
					close(ch)
				}
				delete(processingMap, f)
				processedMap[f] = nil
				mux.Unlock()
				continue
			}
		}(consumerChannel)
	}
}

func pushPacket(cc *ConsumerChannel, result []Packet) {
	outCnt++
	for _, p := range result {
		log.Debugf("[router] pushing converted: %+v", string(p.Body))
		select {
		case cc.Buffer <- p:
		default:
			select {
			case <-cc.Buffer:
				log.Warnf("[router] cache overflowed, popped the oldest message of consumer buffer %v, messages in buffer: %v", cc.UUID, len(cc.Buffer))
				cc.Buffer <- p
			case cc.Buffer <- p:
			}
		}
	}
}
