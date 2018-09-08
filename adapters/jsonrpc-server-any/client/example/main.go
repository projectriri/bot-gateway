package main

import (
	"encoding/json"
	"fmt"
	"github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/client"
	"github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/jsonrpc-any"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/common"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	"github.com/projectriri/bot-gateway/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
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
			data := ubm_api.UBM{}
			json.Unmarshal(pkt.Body, &data)
			// Send message in format telegram-bot-api
			v := url.Values{}
			v.Add("chat_id", data.Message.Chat.CID.ChatID)
			v.Add("text", "Test Send in Telegram-Bot-API Passed!")
			endpoint := fmt.Sprintf(APIEndpoint, "00000000:XXXXXXXXXX_XXXXXXXXXXXXXXXXXXXXXXXX", "sendMessage")
			header := http.Header{}
			header.Set("Content-Type", "application/x-www-form-urlencoded")
			req := common.HTTPRequest{
				Method: "POST",
				URL:    endpoint,
				Header: header,
				Body:   v.Encode(),
			}
			b, _ := json.Marshal(req)
			packet := types.Packet{
				Head: types.Head{
					UUID: utils.GenerateUUID(),
					From: "Test",
					To:   "Telegram",
					Format: types.Format{
						API:      "telegram-bot-api",
						Version:  "latest",
						Protocol: "http",
						Method:   "apirequest",
					},
				},
				Body: b,
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
					Type: "rich_text",
					RichText: &ubm_api.RichText{
						{
							Type: "text",
							Text: "Test Send in UBM-API Passed!",
						},
					},
				},
			}
			b, _ = json.Marshal(ubm)
			packet = types.Packet{
				Head: types.Head{
					UUID: utils.GenerateUUID(),
					From: "Test",
					To:   "Telegram",
					Format: types.Format{
						API:      "ubm-api",
						Version:  "1.0",
						Protocol: "",
						Method:   "send",
					},
				},
				Body: b,
			}
			c.MakeRequest(jsonrpc_any.ChannelProduceRequest{
				c.UUID,
				packet,
			})
		}
	}
}
