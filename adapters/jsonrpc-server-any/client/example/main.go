package main

import (
	"fmt"
	"github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/client"
	"github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/jsonrpc-any"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/ubm-api"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
)

// Telegram constants
const (
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
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
					API:      "ubm-api",
					Version:  "1.0",
					Protocol: "",
					Method:   "receive",
				},
				{
					API:      "telegram-bot-api",
					Version:  "latest",
					Protocol: "http",
					Method:   "update",
				},
			},
		},
	}
	c.Dial()
	pkts, _ := c.GetUpdatesChan(0)
	for pkt := range pkts {
		fmt.Printf("%+v\n", *pkt)
		switch pkt.Head.Format.API {
		case "ubm-api":
			data := pkt.Body.(*ubm_api.UBM)
			// Send message in format telegram-bot-api
			v := url.Values{}
			v.Add("chat_id", data.Message.Chat.CID.ChatID)
			v.Add("text", "Test Send in Telegram-Bot-API Passed!")
			endpoint := fmt.Sprintf(APIEndpoint, "00000000:XXXXXXXXXX_XXXXXXXXXXXXXXXXXXXXXXXX", "sendMessage")
			req, _ := http.NewRequest("POST", endpoint, strings.NewReader(v.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			packet := types.Packet{
				Head: types.Head{
					From: "Test",
					To:   "Telegram",
					Format: types.Format{
						API:      "telegram-bot-api",
						Version:  "latest",
						Protocol: "http",
						Method:   "apirequest",
					},
				},
				Body: req,
			}
			c.MakeRequest(jsonrpc_any.ChannelProduceRequest{
				c.UUID,
				packet,
			})

			// Send Message in format ubm-api
			ubm := ubm_api.UBM{
				Type: "message",
				Message: &ubm_api.Message{
					CID: &ubm_api.CID{
						Messenger: "Telegram",
						ChatID:    data.Message.Chat.CID.ChatID,
					},
					Type: "text",
					RichText: &ubm_api.RichText{
						{
							Text: "Test Send in UBM-API Passed!",
						},
					},
				},
			}
			packet = types.Packet{
				Head: types.Head{
					From: "Test",
					To:   "Telegram",
					Format: types.Format{
						API:      "ubm-api",
						Version:  "1.0",
						Protocol: "",
						Method:   "send",
					},
				},
				Body: &ubm,
			}
			c.MakeRequest(jsonrpc_any.ChannelProduceRequest{
				c.UUID,
				packet,
			})
		}
	}
}
