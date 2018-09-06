package types

type Converter interface {
	BasePlugin
	IsConvertible(from Format, to Format) bool
	Convert(packet Packet, to Format) (bool, []Packet)
}
