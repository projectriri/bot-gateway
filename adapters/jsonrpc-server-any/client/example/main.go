package main

import (
	"fmt"
	"github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/client"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	c := client.Client{}
	c.Init("127.0.0.1:4700", "")
	c.Accept = []router.RoutingRule{
		{
			From: ".*",
			To:   ".*",
			Formats: []types.Format{
				{
					API:      "telegram-bot-api",
					Version:  "latest",
					Protocol: "http",
					Method:   "update",
				},
				{
					API:      "ubm-api",
					Version:  "1.0",
					Protocol: "",
					Method:   "receive",
				},
			},
		},
	}
	c.Dial()
	pkts, _ := c.GetUpdatesChan(0)
	for pkt := range pkts {
		fmt.Println("%+v", *pkt)
	}
}
