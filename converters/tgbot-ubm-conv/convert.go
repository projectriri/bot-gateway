package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"unicode/utf16"
)

func (plugin *Plugin) convertTgUpdateHttpToUbmReceive(packet types.Packet, to types.Format) (bool, []types.Packet) {
	data := packet.Body

	var apiResp tgbotapi.APIResponse
	err := json.Unmarshal(data, &apiResp)
	if err != nil {
		return false, nil
	}
	var updates []tgbotapi.Update
	json.Unmarshal(apiResp.Result, &updates)

	result := make([]types.Packet, 0)

	for _, update := range updates {
		if update.EditedMessage != nil || update.Message != nil || update.ChannelPost != nil || update.EditedChannelPost != nil {
			if update.EditedMessage != nil {
				update.Message = update.EditedMessage
			} else if update.EditedChannelPost != nil {
				update.Message = update.EditedChannelPost
			} else if update.ChannelPost != nil {
				update.Message = update.ChannelPost
			}
			self := plugin.getMe(packet.Head.From)
			ubm := ubm_api.UBM{
				Type: "message",
				Date: int64(update.Message.Date),
				Self: self,
				Message: &ubm_api.Message{
					ID: strconv.Itoa(update.Message.MessageID),
				},
			}
			if update.ChannelPost != nil || update.EditedChannelPost != nil {
				ubm.Message.From = &ubm_api.User{
					DisplayName: update.Message.Chat.Title,
					FirstName:   update.Message.Chat.Title,
					UID: ubm_api.UID{
						Messenger: packet.Head.From,
						ID:        strconv.FormatInt(update.Message.Chat.ID, 10),
						Username:  update.Message.Chat.UserName,
					},
					PrivateChat: ubm_api.CID{
						Messenger: packet.Head.From,
						ChatID:    strconv.FormatInt(update.Message.Chat.ID, 10),
						ChatType:  "channel",
					},
				}
			} else {
				ubm.Message.From = &ubm_api.User{
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
				}
			}
			ubm.Message.Chat = &ubm_api.Chat{
				Title:       update.Message.Chat.Title,
				Description: update.Message.Chat.Description,
				CID: ubm_api.CID{
					Messenger: packet.Head.From,
					ChatID:    strconv.FormatInt(update.Message.Chat.ID, 10),
					ChatType:  getTelegramChatType(update.Message.Chat),
				},
			}
			if update.Message.ReplyToMessage != nil {
				if self != nil && strconv.Itoa(update.Message.ReplyToMessage.From.ID) == self.UID.ID {
					ubm.Message.IsMessageToMe = true
				}
			}
			if update.Message.Sticker != nil {
				ubm.Message.Type = "sticker"
				ubm.Message.Sticker = &ubm_api.Sticker{
					Messenger: packet.Head.From,
					ID:        update.Message.Sticker.FileID,
					Image: &ubm_api.Image{
						Messenger: packet.Head.From,
						Width:     update.Message.Sticker.Width,
						Height:    update.Message.Sticker.Height,
						FileID:    update.Message.Sticker.FileID,
						URL:       plugin.getFileURL(update.Message.Sticker.FileID, packet.Head.From),
						FileSize:  update.Message.Sticker.FileSize,
					},
				}
			} else if update.Message.Location != nil {
				ubm.Message.Type = "location"
				ubm.Message.Location = &ubm_api.Location{
					Latitude:  update.Message.Location.Latitude,
					Longitude: update.Message.Location.Longitude,
				}
			} else if update.Message.Audio != nil {
				ubm.Message.Type = "voice"
				ubm.Message.Voice = &ubm_api.Voice{
					Messenger: packet.Head.From,
					Format:    update.Message.Audio.MimeType,
					Duration:  update.Message.Audio.Duration,
					FileID:    update.Message.Audio.FileID,
					URL:       plugin.getFileURL(update.Message.Audio.FileID, packet.Head.From),
					FileSize:  update.Message.Audio.FileSize,
				}
			} else if update.Message.Voice != nil {
				ubm.Message.Type = "voice"
				ubm.Message.Voice = &ubm_api.Voice{
					Messenger: packet.Head.From,
					Format:    update.Message.Voice.MimeType,
					Duration:  update.Message.Voice.Duration,
					FileID:    update.Message.Voice.FileID,
					URL:       plugin.getFileURL(update.Message.Voice.FileID, packet.Head.From),
					FileSize:  update.Message.Voice.FileSize,
				}
			} else if update.Message.Photo != nil && len(*update.Message.Photo) > 0 {
				ubm.Message.Type = "rich_text"
				richText := make(ubm_api.RichText, 0)
				photo := (*update.Message.Photo)[len(*update.Message.Photo)-1]
				richText = append(richText, ubm_api.RichTextElement{
					Type: "image",
					Image: &ubm_api.Image{
						Messenger: packet.Head.From,
						Width:     photo.Width,
						Height:    photo.Height,
						FileID:    photo.FileID,
						URL:       plugin.getFileURL(photo.FileID, packet.Head.From),
						FileSize:  photo.FileSize,
					},
				})
				if update.Message.Caption != "" {
					if self != nil && strings.Contains(update.Message.Caption, "@"+self.UID.Username) {
						ubm.Message.IsMessageToMe = true
					}
					richText = append(richText, ubm_api.RichTextElement{
						Type: "text",
						Text: update.Message.Caption,
					})
				}
				ubm.Message.RichText = &richText
			} else if update.Message.Text != "" {
				r := []rune(update.Message.Text)
				u16 := utf16.Encode(r)
				cur := 0
				tmp := ""
				ubm.Message.Type = "rich_text"
				ubm.Message.RichText = &ubm_api.RichText{}
				push := func() {
					if len(tmp) > 0 {
						*ubm.Message.RichText = append(*ubm.Message.RichText, ubm_api.RichTextElement{
							Type: "text",
							Text: tmp,
						})
						tmp = ""
					}
				}
				if update.Message.Entities != nil {
					for _, entity := range *update.Message.Entities {
						cat := func() {
							if entity.Offset > cur {
								tmp += string(utf16.Decode(u16[cur:entity.Offset]))
							}
						}
						switch entity.Type {
						case "bot_command":
							cat()
							botCommand := string(utf16.Decode(u16[entity.Offset : entity.Offset+entity.Length]))
							if self != nil && len(botCommand) > len("@"+self.UID.Username) &&
								botCommand[len(botCommand)-len("@"+self.UID.Username):] == "@"+self.UID.Username {
								botCommand = botCommand[:len(botCommand)-len("@"+self.UID.Username)]
								ubm.Message.IsMessageToMe = true
							}
							tmp += botCommand
							cur = entity.Offset + entity.Length
						case "mention":
							cat()
							push()
							mention := string(utf16.Decode(u16[entity.Offset : entity.Offset+entity.Length]))
							*ubm.Message.RichText = append(*ubm.Message.RichText, ubm_api.RichTextElement{
								Type: "at",
								At: &ubm_api.At{
									DisplayName: mention,
									UID: ubm_api.UID{
										Messenger: packet.Head.From,
										Username:  mention[1:],
									},
								},
							})
							if self != nil && mention == "@"+self.UID.Username {
								ubm.Message.IsMessageToMe = true
							}
							cur = entity.Offset + entity.Length
						case "text_mention":
							cat()
							push()
							mention := string(utf16.Decode(u16[entity.Offset : entity.Offset+entity.Length]))
							*ubm.Message.RichText = append(*ubm.Message.RichText, ubm_api.RichTextElement{
								Type: "at",
								At: &ubm_api.At{
									DisplayName: mention,
									UID: ubm_api.UID{
										Messenger: packet.Head.From,
										ID:        strconv.Itoa(entity.User.ID),
									},
								},
							})
							if self != nil && strconv.Itoa(entity.User.ID) == self.UID.ID {
								ubm.Message.IsMessageToMe = true
							}
							cur = entity.Offset + entity.Length
						}
					}
				}
				if len(update.Message.Text) > cur {
					tmp += string(utf16.Decode(u16[cur:]))
				}
				push()
			} else {
				// notice or something we cannot handle
				b, _ := json.Marshal(update.Message)
				log.Warnf("[tgbot-ubm-conv] cannot convert message from tgbot to ubm: "+
					"%s", string(b))
				continue
			}
			b, _ := json.Marshal(ubm)
			p := types.Packet{
				Head: packet.Head,
				Body: b,
			}
			p.Head.Format = to
			p.Head.Format.Version = UBMAPIVersion
			result = append(result, p)
		}
	}

	return len(result) > 0, result
}

func (plugin *Plugin) convertUbmSendToTgApiRequestHttp(packet types.Packet, to types.Format) (bool, []types.Packet) {
	data := ubm_api.UBM{}
	err := json.Unmarshal(packet.Body, &data)
	if err != nil {
		return false, nil
	}
	p := types.Packet{
		Head: packet.Head,
	}
	p.Head.Format = to
	p.Head.Format.Version = TelegramBotAPIVersion

	result := make([]types.Packet, 0)

	switch data.Type {
	case "message":
		if data.Message == nil {
			return false, nil
		}
		v := make(map[string]string)
		v["chat_id"] = data.Message.CID.ChatID
		if data.Message.ReplyID != "" {
			v["reply_to_message_id"] = data.Message.ReplyID
		}
		if data.Message.EditID != "" && data.Message.RichText != nil {
			v["message_id"] = data.Message.EditID
			var text string
			for _, elem := range *data.Message.RichText {
				text += elem.Text
			}
			v["text"] = text
			p.Body, _ = json.Marshal(newMessageRequest("editMessageText", v))
			result = append(result, p)
			break
		}
		if data.Message.DeleteID != "" {
			v["message_id"] = data.Message.DeleteID
			p.Body, _ = json.Marshal(newMessageRequest("deleteMessage", v))
			result = append(result, p)
			break
		}
		switch data.Message.Type {
		case "voice":
			if data.Message.Voice == nil {
				return false, nil
			}
			if data.Message.Voice.URL != "" {
				v["voice"] = data.Message.Voice.URL
				p.Body, _ = json.Marshal(newMessageRequest("sendVoice", v))
				result = append(result, p)
				break
			}
			if data.Message.Voice.Data != nil {
				p.Body, _ = json.Marshal(newFileRequest("sendVoice", v, map[string][]byte{
					"voice": data.Message.Voice.Data,
				}))
				result = append(result, p)
				break
			}
			return false, nil
		case "location":
			if data.Message.Location == nil {
				return false, nil
			}
			v["latitude"] = strconv.FormatFloat(data.Message.Location.Latitude, 'f', 6, 64)
			v["longitude"] = strconv.FormatFloat(data.Message.Location.Longitude, 'f', 6, 64)
			p.Body, _ = json.Marshal(newMessageRequest("sendLocation", v))
			result = append(result, p)
			break
		case "sticker":
			if data.Message.Sticker == nil {
				return false, nil
			}
			if data.Message.Sticker.ID != "" {
				v["sticker"] = data.Message.Sticker.ID
				p.Body, _ = json.Marshal(newMessageRequest("sendSticker", v))
				result = append(result, p)
				break
			}
			if data.Message.Sticker.Image != nil {
				if data.Message.Sticker.Image.URL != "" {
					v["sticker"] = data.Message.Sticker.Image.URL
					p.Body, _ = json.Marshal(newMessageRequest("sendSticker", v))
					result = append(result, p)
					break
				}
				if data.Message.Sticker.Image.Data != nil {
					p.Body, _ = json.Marshal(newFileRequest("sendSticker", v, map[string][]byte{
						"sticker": data.Message.Sticker.Image.Data,
					}))
					result = append(result, p)
					break
				}
				return false, nil
			}
			return false, nil
		case "rich_text":
			v2 := mapCopy(v)

			// format ats and concat neighbor texts
			txtTmp := make([]string, 0)
			rtxArr := make(ubm_api.RichText, 0)
			for _, elem := range *data.Message.RichText {
				switch elem.Type {
				case "image":
					if len(txtTmp) != 0 {
						rtxArr = append(rtxArr, ubm_api.RichTextElement{
							Type: "text",
							Text: strings.Join(txtTmp, " "),
						})
					}
					txtTmp = make([]string, 0)
					rtxArr = append(rtxArr, elem)
				case "text":
					t := plainToMarkdown(elem.Text)
					if t != "" {
						txtTmp = append(txtTmp, t)
					}
				case "at":
					if elem.At == nil {
						continue
					}
					if elem.At.UID.Username != "" {
						txtTmp = append(txtTmp, fmt.Sprintf("@%s", plainToMarkdown(elem.At.UID.Username)))
					} else {
						txtTmp = append(txtTmp, fmt.Sprintf("[%s](tg://user?id=%s)", elem.At.DisplayName, elem.At.UID.ID))
					}
				}
			}
			if len(txtTmp) != 0 {
				rtxArr = append(rtxArr, ubm_api.RichTextElement{
					Type: "text",
					Text: strings.Join(txtTmp, " "),
				})
			}

			photos := make(map[string][]byte)
			photoParams := make([]PhotoConfig, 0)
			for _, elem := range rtxArr {
				switch elem.Type {
				case "text":
					if len(photoParams) == 1 {
						// send one image with caption
						v["caption"] = elem.Text
						v["parse_mode"] = "Markdown"
						if len(photos) == 0 {
							v["photo"] = photoParams[0].Media
							p.Body, _ = json.Marshal(newMessageRequest("sendPhoto", v))
						} else {
							p.Body, _ = json.Marshal(newFileRequest("sendPhoto", v, photos))
						}
						result = append(result, p)
						v = mapCopy(v2)
						photos = make(map[string][]byte)
						photoParams = make([]PhotoConfig, 0)
					} else if len(photoParams) == 0 {
						// send a text message
						v["text"] = elem.Text
						v["parse_mode"] = "Markdown"
						p.Body, _ = json.Marshal(newMessageRequest("sendMessage", v))
						result = append(result, p)
						v = mapCopy(v2)
					} else {
						// send images
						b, _ := json.Marshal(photoParams)
						v["media"] = string(b)
						p.Body, _ = json.Marshal(newFileRequest("sendmediagroup", v, photos))
						result = append(result, p)
						v = mapCopy(v2)
						photos = make(map[string][]byte)
						photoParams = make([]PhotoConfig, 0)
						// send text message
						v["text"] = elem.Text
						v["parse_mode"] = "Markdown"
						p.Body, _ = json.Marshal(newMessageRequest("sendMessage", v))
						result = append(result, p)
						v = mapCopy(v2)
					}
				case "image":
					if elem.Image == nil {
						continue
					}
					var field string
					if len(photoParams) == 0 {
						field = "photo"
					} else {
						field = fmt.Sprintf("photo%d", len(photoParams))
					}
					if elem.Image.FileID != "" {
						photoParams = append(photoParams, PhotoConfig{
							Type:  "photo",
							Media: elem.Image.FileID,
						})
						break
					}
					if elem.Image.URL != "" {
						photoParams = append(photoParams, PhotoConfig{
							Type:  "photo",
							Media: elem.Image.URL,
						})
						break
					}
					if elem.Image.Data != nil {
						photoParams = append(photoParams, PhotoConfig{
							Type:  "photo",
							Media: fmt.Sprintf("attach://%s", field),
						})
						photos[field] = elem.Image.Data
					}
				}
			}
			if len(photoParams) == 1 {
				if len(photos) == 0 {
					v["photo"] = photoParams[0].Media
					p.Body, _ = json.Marshal(newMessageRequest("sendPhoto", v))
				} else {
					p.Body, _ = json.Marshal(newFileRequest("sendPhoto", v, photos))
				}
				result = append(result, p)
			} else if len(photoParams) > 1 {
				b, _ := json.Marshal(photoParams)
				v["media"] = string(b)
				p.Body, _ = json.Marshal(newFileRequest("sendmediagroup", v, photos))
				result = append(result, p)
			}
			break
		}

	}
	return len(result) > 0, result
}
