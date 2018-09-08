package main

import (
	"fmt"
	"github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/client"
)

func main() {
	c := client.Client{}
	c.Init("127.0.0.1:4700", "")
	pkts, _ := c.GetUpdatesChan(0)
	for pkt := range pkts {
		fmt.Println("%+v", *pkt)
	}
}
