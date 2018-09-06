package router

import (
	. "github.com/projectriri/bot-gateway/types"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func route() {
	for {
		pkt := <-producerBuffer
		log.Debugf("[router] pkt: %+v", pkt.Head)

		from := strings.ToLower(pkt.Head.From)
		to := strings.ToLower(pkt.Head.To)

		for _, pcc := range consumerChannelPool {
			go func(cc *ConsumerChannel) {
				log.Debugf("[router] pkt: %v TRYING cc: %+v", pkt.Head.UUID, cc)
				var formats []Format
				for _, ac := range cc.Accept {
					f, _ := regexp.MatchString(ac.From, from)
					t, _ := regexp.MatchString(ac.To, to)
					if f && t {
						formats = ac.Formats
						break
					}
				}

				if formats == nil {
					return
				}

				log.Debugf("[router] pkt: %v ACCEPTED BY cc: %+v", pkt.Head.UUID, cc)

				for _, format := range formats {

					if strings.ToLower(pkt.Head.Format.API) == strings.ToLower(format.API) &&
						strings.ToLower(pkt.Head.Format.Method) == strings.ToLower(format.Method) &&
						strings.ToLower(pkt.Head.Format.Protocol) == strings.ToLower(format.Protocol) {
						cc.Buffer <- pkt
						return
					}

					for _, cvt := range converters {
						if cvt.IsConvertible(pkt.Head.Format, format) {
							ok, result := cvt.Convert(pkt, format)
							if ok && result != nil {
								for _, p := range result {
									cc.Buffer <- p
								}
								return
							}
						}
					}
				}
			}(pcc)
		}
	}
}
