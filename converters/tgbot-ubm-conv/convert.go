package main

import (
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/ubm-api"
	"net/url"
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
				ubm.Message.Type = "location"
				ubm.Message.Location = &ubm_api.Location{
					Latitude:  update.Message.Location.Latitude,
					Longitude: update.Message.Location.Longitude,
				}
			} else if update.Message.Audio != nil {
				// TODO
				return false
			} else if update.Message.Photo != nil {
				ubm.Message.Type = "rich_text"
				// TODO
				// for _, photo := range *update.Message.Photo {
				//
				// }
				ubm.Message.RichText = &ubm_api.RichText{
					{
						Type: "text",
						Text: update.Message.Text,
					},
				}
			} else if update.Message.Text != "" {
				ubm.Message.Type = "rich_text"
				ubm.Message.RichText = &ubm_api.RichText{
					{
						Type: "text",
						Text: update.Message.Text,
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

func convertUbmSendToTgApiRequestHttp(packet router.Packet, to router.Format, ch router.Buffer) bool {
	data, ok := packet.Body.(ubm_api.UBM)
	if !ok {
		return false
	}
	p := router.Packet{
		Head: packet.Head,
	}
	p.Head.Format = to
	switch data.Type {
	case "message":
		if data.Message == nil {
			return false
		}
		v := url.Values{}
		v.Add("chat_id", data.Message.CID.ChatID)
		if data.Message.ReplyID != "" {
			v.Add("reply_to_message_id", data.Message.ReplyID)
		}
		if data.Message.ForwardID != "" {
			v.Add("from_chat_id", data.Message.ForwardFromChat.CID.ChatID)
			v.Add("message_id", data.Message.ForwardID)
			p.Body = newMessageRequest("forwardMessage", v)
			ch <- p
			return true
		}
		switch data.Message.Type {
		case "audio":
			// TODO
			return false
		case "location":
			v.Add("latitude", strconv.FormatFloat(data.Message.Location.Latitude, 'f', 6, 64))
			v.Add("longitude", strconv.FormatFloat(data.Message.Location.Longitude, 'f', 6, 64))
			p.Body = newMessageRequest("sendLocation", v)
			ch <- p
			return true
		case "sticker":
			if data.Message.Sticker.ID != "" {
				v.Add("sticker", data.Message.Sticker.ID)
				p.Body = newMessageRequest("sendSticker", v)
				ch <- p
				return true
			} else {
				// TODO
				return false
			}
		case "rich_text":
			photoTmp := false
			for _, elem := range *data.Message.RichText {
				switch elem.Type {
				case "styled_text":
					fallthrough
				case "text":
					if !photoTmp {
						if elem.Type == "styled_text" {
							v.Add("text", elem.StyledText.Text)
							v.Add("parse_mode", elem.StyledText.Format)
						} else {
							v.Add("text", elem.Text)
						}
						p.Body = newMessageRequest("sendMessage", v)
						ch <- p
					} else {
						if elem.Type == "styled_text" {
							v.Add("caption", elem.StyledText.Text)
						} else {
							v.Add("caption", elem.Text)
						}
						p.Body = newMessageRequest("sendPhoto", v)
						ch <- p
						photoTmp = false
					}
				case "image":
					// TODO
				}
			}
			if photoTmp {
				// TODO
			}
			return true
		}

	}
	return false
}
