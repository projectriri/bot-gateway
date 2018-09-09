package types

type Converter interface {
	Adapter
	IsConvertible(from Format, to Format) bool
	Convert(packet Packet, to Format) (bool, []Packet)
}
