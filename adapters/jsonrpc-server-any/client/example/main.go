package main

import (
	"fmt"
	"github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/client"
)

func main() {
	client.Init("127.0.0.1:4700", "")
	pkts, _ := client.GetUpdatesChan(0)
	for pkt := range pkts {
		fmt.Println("%+v", *pkt)
	}
}
