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
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	// in an actual fly, place a legal uuid here
	uuid := ""
	c.Init("127.0.0.1:4700", uuid)
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
			// keep the token as the placeholder, this will be handled by the Telegram adapter
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

			// Send message in format ubm-api

			file, err := os.Open("projectriri.jpg")
			if err != nil {
				fmt.Println(err)
			}
			fileContents, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
			}

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
						{
							Type: "at",
							At: &ubm_api.At{
								DisplayName: "梨子",
								UID: ubm_api.UID{
									Messenger: "Telegram",
									ID:        "8964",
								},
							},
						},
						{
							Type: "text",
							Text: "Text 2: There is an at and an image before this text.",
						},
						{
							Type: "text",
							Text: "Text 2: There is an image before this text.",
						},
						{
							Type: "image",
							Image: &ubm_api.Image{
								Data: &fileContents,
							},
						},
						{
							Type: "image",
							Image: &ubm_api.Image{
								Data: &fileContents,
							},
						},
						{
							Type: "at",
							At: &ubm_api.At{
								DisplayName: "梨子",
								UID: ubm_api.UID{
									Messenger: "Telegram",
									ID:        "8964",
									Username:  "example_user",
								},
							},
						},
						{
							Type: "text",
							Text: "Text 3: There is an at and 2 images before this text.",
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
