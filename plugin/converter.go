package plugin

import "github.com/projectriri/bot-gateway/router"

type Converter interface {
	BasePlugin
	IsConvertible(from router.Format, to router.Format) bool
	Convert(packet router.Packet, to router.Format, destination router.Buffer) bool
}
