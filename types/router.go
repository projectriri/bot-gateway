package types

type Packet struct {
	Head Head
	Body interface{}
}

type Head struct {
	UUID                   string
	From                   string
	To                     string
	ReplyToUUID            string
	AcknowledgeChannelUUID string
	Format                 Format
}

type Format struct {
	API      string
	Version  string
	Method   string
	Protocol string
}

type Buffer chan Packet
