package main

import (
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/ubm-api"
	"strconv"
)

func convertTgUpdateHttpToUbmReceive(packet router.Packet, to router.Format, ch router.Buffer) bool {
	data, ok := packet.Body.([]byte)
	if !ok {
		return false
	}
	var apiResp tgbotapi.APIResponse
	err := json.Unmarshal(data, apiResp)
	if err != nil {
		return false
	}
	var updates []tgbotapi.Update
	json.Unmarshal(apiResp.Result, &updates)
	flag := false
	for _, update := range updates {
		if update.EditedMessage != nil || update.Message != nil || update.ChannelPost != nil || update.EditedChannelPost != nil {
			if update.EditedMessage != nil {
				update.Message = update.EditedMessage
			} else if update.EditedChannelPost != nil {
				update.Message = update.EditedChannelPost
			} else if update.EditedMessage != nil {
				update.Message = update.ChannelPost
			}
			ubm := ubm_api.UBM{
				Type: "message",
				Message: &ubm_api.Message{
					ID: strconv.Itoa(update.Message.MessageID),
					From: &ubm_api.User{
						DisplayName: update.Message.From.FirstName,
						FirstName:   update.Message.From.FirstName,
						LastName:    update.Message.From.LastName,
						UID: ubm_api.UID{
							Messenger: packet.Head.From,
							ID:        strconv.Itoa(update.Message.From.ID),
							Username:  update.Message.From.UserName,
						},
						PrivateChat: ubm_api.CID{
							Messenger: packet.Head.From,
							ChatID:    strconv.Itoa(update.Message.From.ID),
							ChatType:  "private",
						},
					},
					Chat: &ubm_api.Chat{
						Title:       update.Message.Chat.Title,
						Description: update.Message.Chat.Description,
						CID: ubm_api.CID{
							Messenger: packet.Head.From,
							ChatID:    strconv.FormatInt(update.Message.Chat.ID, 10),
							ChatType:  getTelegramChatType(update.Message.Chat),
						},
					},
				},
			}
			if update.Message.Sticker != nil {
				ubm.Message.Type = "sticker"
				ubm.Message.Sticker = &ubm_api.Sticker{
					ID: update.Message.Sticker.FileID,
				}
			} else if update.Message.Location != nil {
				ubm.Message.Location = &ubm_api.Location{
					Latitude:  update.Message.Location.Latitude,
					Longitude: update.Message.Location.Longitude,
				}
			} else if update.Message.Photo != nil {
				// for _, photo := range *update.Message.Photo {
				//
				// }
				ubm.Message.RichText = &ubm_api.RichText{
					{
						Type: "text",
						Data: update.Message.Caption,
					},
				}
			} else if update.Message.Text != "" {
				ubm.Message.RichText = &ubm_api.RichText{
					{
						Type: "text",
						Data: update.Message.Text,
					},
				}
			}
			p := router.Packet{
				Head: packet.Head,
				Body: ubm,
			}
			p.Head.Format = to
			ch <- p
			flag = true
		}
	}
	return flag
}
