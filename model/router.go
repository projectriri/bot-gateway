package model

type Packet struct {
	Head Head
	Body []byte
}

type Head struct {
	UUID      string
	From      string
	To        string
	ReplyUUID string
	Level     int
	Format    Format
}

type Format struct {
	API      string
	Version  string
	Protocol string
}

type Buffer chan Packet

type Channel struct {
	UUID   string
	Buffer *Buffer
}

type ProducerChannel struct {
	Channel
	AcknowlegeBuffer *Buffer
}

type ConsumerChannel struct {
	Channel
	Accept []RoutingRule
}

type RoutingRule struct {
	From    string
	To      string
	Level   int
	Formats []Format
}
