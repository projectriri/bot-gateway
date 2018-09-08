package types

import "encoding/json"

type Packet struct {
	Head Head            `json:"head"`
	Body json.RawMessage `json:"body"`
}

type Head struct {
	UUID                   string `json:"uuid"`
	From                   string `json:"from"`
	To                     string `json:"to"`
	ReplyToUUID            string `json:"reply_to_uuid"`
	AcknowledgeChannelUUID string `json:"acknowledge_channel_uuid"`
	Format                 Format `json:"format"`
}

type Format struct {
	API      string `json:"api"`
	Version  string `json:"version"`
	Method   string `json:"method"`
	Protocol string `json:"protocol"`
}

type Buffer chan Packet
