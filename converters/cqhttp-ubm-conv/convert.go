package main

import (
	"encoding/json"
	"fmt"
	"github.com/catsworld/qq-bot-api"
	"github.com/catsworld/qq-bot-api/cqcode"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strconv"
)

func (p *Plugin) convertQQEventWSToUbmReceive(packet types.Packet, to types.Format) (bool, []types.Packet) {
	// parse cqhttp update
	update := qqbotapi.Update{}
	err := json.Unmarshal(packet.Body, &update)
	if err != nil {
		log.Errorf("[cqhttp-ubm-conv] failed to unmarshal json (%v)", err)
		return false, nil
	}
	update.ParseRawMessage()

	// only convert message now
	// TODO: convert notice and request
	if update.PostType != "message" {
		return false, nil
	}

	if update.Sender == nil {
		// fill in cq update.message.from
		switch update.Message.Chat.Type {
		case "group":
			res := p.makeRequest(
				packet.Head.From,
				"get_group_member_info",
				map[string]interface{}{
					"group_id": update.Message.Chat.ID,
					"user_id":  update.Message.From.ID,
				},
			)
			if res != nil {
				user := qqbotapi.User{}
				if err := json.Unmarshal(res, &user); err == nil {
					update.Message.From = &user
				}
			}
		default:
			res := p.makeRequest(
				packet.Head.From,
				"get_stranger_info",
				map[string]interface{}{
					"user_id": update.Message.From.ID,
				},
			)
			if res != nil {
				user := qqbotapi.User{}
				if err := json.Unmarshal(res, &user); err == nil {
					update.Message.From = &user
				}
			}
		}
	}

	// convert to UBM
	ubm := ubm_api.UBM{}
	ubm.Type = update.PostType
	ubm.Date = update.Time
	ubm.Self = p.getMe(packet.Head.From)
	ubm.Message = &ubm_api.Message{
		ID: strconv.FormatInt(update.Message.MessageID, 10),
		From: &ubm_api.User{
			DisplayName: update.Message.From.Name(),
			FirstName:   update.Message.From.NickName,
			LastName:    "",
			UID: ubm_api.UID{
				Messenger: packet.Head.From,
				ID:        strconv.FormatInt(update.Message.From.ID, 10),
				Username:  "",
			},
			PrivateChat: ubm_api.CID{
				Messenger: packet.Head.From,
				ChatID:    strconv.FormatInt(update.Message.From.ID, 10),
				ChatType:  "private",
			},
		},
		Chat: &ubm_api.Chat{
			CID: ubm_api.CID{
				Messenger: packet.Head.From,
				ChatID:    strconv.FormatInt(update.Message.Chat.ID, 10),
				ChatType:  update.Message.Chat.Type,
			},
		},
	}

	ubmHandleFinishFlag := false
	// handle voice and location
	if len(*update.Message.Message) == 1 {
		mediaInterface := (*update.Message.Message)[0]
		switch media := mediaInterface.(type) {
		case *cqcode.Record:
			ubm.Message.Type = "voice"
			ubm.Message.Voice = &ubm_api.Voice{
				Messenger: packet.Head.From,
				URL:       media.URL,
				FileID:    media.FileID,
			}
			ubmHandleFinishFlag = true
		case *cqcode.Location:
			ubm.Message.Type = "location"
			ubm.Message.Location = &ubm_api.Location{
				Content:   media.Content,
				Latitude:  media.Latitude,
				Longitude: media.Longitude,
				Title:     media.Title,
			}
			ubmHandleFinishFlag = true
		}
	}
	// handle rich text
	if !ubmHandleFinishFlag {
		richTexts := make(ubm_api.RichText, 0)
		for _, mediaInterface := range *update.Message.Message {
			switch media := mediaInterface.(type) {
			case *cqcode.Text:
				richTexts = append(richTexts, ubm_api.RichTextElement{
					Type: "text",
					Text: media.Text,
				})
			case *cqcode.At:
				richTexts = append(richTexts, ubm_api.RichTextElement{
					Type: "at",
					At: &ubm_api.At{
						UID: ubm_api.UID{
							Messenger: packet.Head.From,
							ID:        media.QQ,
						},
					},
				})
				// check is_message_to_me
				if ubm.Self != nil && (media.QQ == "all" || media.QQ == ubm.Self.UID.ID) {
					ubm.Message.IsMessageToMe = true
				}
			case *cqcode.Image:
				richTexts = append(richTexts, ubm_api.RichTextElement{
					Type: "image",
					Image: &ubm_api.Image{
						Messenger: packet.Head.From,
						FileID:    media.FileID,
						URL:       media.URL,
					},
				})
			}
		}
		if len(richTexts) > 0 {
			ubm.Message.Type = "rich_text"
			ubm.Message.RichText = &richTexts
			ubmHandleFinishFlag = true
		}
	}
	// not voice, location nor rich text
	if !ubmHandleFinishFlag {
		return false, nil
	}

	// put UBM into packets
	packets := make([]types.Packet, 1)
	packets[0].Head = packet.Head
	packets[0].Head.Format = to
	packets[0].Head.Format.Version = UBMAPIVersion
	packets[0].Body, _ = json.Marshal(ubm)

	return true, packets
}

func (p *Plugin) convertUbmVoiceToCQRecord(voice *ubm_api.Voice) (bool, cqcode.Media) {
	if voice.Data != nil && len(voice.Data) > 0 {
		if rc, err := qqbotapi.NewRecordBase64(voice.Data); err != nil {
			log.Errorf("[cqhttp-ubm-conv] failed to encode voice data to base64 (%v)", err)
			return false, nil
		} else {
			return true, rc
		}
	} else if voice.URL != "" {
		if u, err := url.Parse(voice.URL); err != nil {
			log.Errorf("[cqhttp-ubm-conv] failed to parse voice url (%v)", err)
			return false, nil
		} else {
			rc := qqbotapi.NewRecordWeb(u)
			rc.DisableCache()
			return true, rc
		}
	} else if voice.FileID != "" {
		return true, &cqcode.Record{
			FileID: voice.FileID,
		}
	} else {
		log.Errorf("[cqhttp-ubm-conv] unknown voice format")
		return false, nil
	}
}

func (p *Plugin) convertUbmImageToCQImage(image *ubm_api.Image) (bool, cqcode.Media) {
	if image.Data != nil && len(image.Data) > 0 {
		if img, err := qqbotapi.NewImageBase64(image.Data); err != nil {
			log.Errorf("[cqhttp-ubm-conv] failed to encode image data to base64 (%v)", err)
			return false, nil
		} else {
			return true, img
		}
	} else if image.URL != "" {
		if u, err := url.Parse(image.URL); err != nil {
			log.Errorf("[cqhttp-ubm-conv] failed to parse image url (%v)", err)
			return false, nil
		} else {
			img := qqbotapi.NewImageWeb(u)
			img.DisableCache()
			return true, img
		}
	} else if image.FileID != "" {
		return true, &cqcode.Image{
			FileID: image.FileID,
		}
	} else {
		log.Errorf("[cqhttp-ubm-conv] unknown image format")
		return false, nil
	}
}

func (p *Plugin) convertUbmSendToQQApiRequestWS(packet types.Packet, to types.Format) (bool, []types.Packet) {
	var data []byte

	// unmarshal UBM
	ubm := ubm_api.UBM{}
	if err := json.Unmarshal(packet.Body, &ubm); err != nil {
		log.Errorf("[cqhttp-ubm-conv] failed to unmarshal json (%v)", err)
		return false, nil
	}

	// TODO: handle UBM action
	if ubm.Type != "message" || ubm.Message == nil {
		return false, nil
	}

	// handle UBM message accordingly
	if ubm.Message.DeleteID != "" {
		// delete message
		data = []byte(fmt.Sprintf(
			`{"action": "delete_msg", "params": {"message_id": %s}}`,
			ubm.Message.DeleteID,
		))
	} else {
		// else send message
		message := cqcode.Message{}
		switch ubm.Message.Type {
		case "location":
			if ubm.Message.Location == nil {
				log.Errorf("[cqhttp-ubm-conv] ubm.Location is nil")
				return false, nil
			}
			message.Append(&cqcode.Location{
				Content:   ubm.Message.Location.Content,
				Latitude:  ubm.Message.Location.Latitude,
				Longitude: ubm.Message.Location.Longitude,
				Title:     ubm.Message.Location.Title,
			})
		case "voice":
			if ubm.Message.Voice == nil {
				log.Errorf("[cqhttp-ubm-conv] ubm.Voice is nil")
				return false, nil
			}
			if ok, rc := p.convertUbmVoiceToCQRecord(ubm.Message.Voice); ok {
				message.Append(rc)
			} else {
				return false, nil
			}
		case "sticker":
			if ubm.Message.Sticker == nil {
				log.Errorf("[cqhttp-ubm-conv] ubm.Sticker is nil")
				return false, nil
			}
			if ubm.Message.Sticker.PackID != "" {
				switch ubm.Message.Sticker.PackID {
				case "face":
					if fid, err := strconv.Atoi(ubm.Message.Sticker.ID); err != nil {
						log.Errorf("[cqhttp-ubm-conv] failed to parse face id")
						return false, nil
					} else {
						message.Append(&cqcode.Face{
							FaceID: fid,
						})
					}
				case "bface":
					if fid, err := strconv.Atoi(ubm.Message.Sticker.ID); err != nil {
						log.Errorf("[cqhttp-ubm-conv] failed to parse Bface id")
						return false, nil
					} else {
						message.Append(&cqcode.Bface{
							BfaceID: fid,
						})
					}
				case "sface":
					if fid, err := strconv.Atoi(ubm.Message.Sticker.ID); err != nil {
						log.Errorf("[cqhttp-ubm-conv] failed to parse Sface id")
						return false, nil
					} else {
						message.Append(&cqcode.Sface{
							SfaceID: fid,
						})
					}
				case "rps":
					if fid, err := strconv.Atoi(ubm.Message.Sticker.ID); err != nil {
						log.Errorf("[cqhttp-ubm-conv] failed to parse rps id")
						return false, nil
					} else {
						message.Append(&cqcode.Rps{
							Type: fid,
						})
					}
				case "dice":
					if fid, err := strconv.Atoi(ubm.Message.Sticker.ID); err != nil {
						log.Errorf("[cqhttp-ubm-conv] failed to parse dice id")
						return false, nil
					} else {
						message.Append(&cqcode.Dice{
							Type: fid,
						})
					}
				case "shake":
					message.Append(&cqcode.Shake{})
				default:
					log.Errorf("[cqhttp-ubm-conv] unknown sticker pack id")
					return false, nil
				}
			} else if ubm.Message.Sticker.Image != nil {
				if ok, img := p.convertUbmImageToCQImage(ubm.Message.Sticker.Image); ok {
					message.Append(img)
				} else {
					return false, nil
				}
			} else {
				log.Errorf("[cqhttp-ubm-conv] unknown sticker format")
				return false, nil
			}
		case "rich_text":
			if ubm.Message.RichText == nil {
				log.Errorf("[cqhttp-ubm-conv] ubm.RichText is nil")
				return false, nil
			}
			for _, richTextElem := range *ubm.Message.RichText {
				switch richTextElem.Type {
				case "text":
					message.Append(&cqcode.Text{
						Text: richTextElem.Text,
					})
				case "at":
					message.Append(&cqcode.At{
						QQ: richTextElem.At.UID.ID,
					})
				case "image":
					if ok, img := p.convertUbmImageToCQImage(richTextElem.Image); ok {
						message.Append(img)
					} else {
						continue
					}
				}
			}
		}
		if len(message) == 0 {
			log.Warnf("[cqhttp-ubm-conv] converted cqhttp message is empty, ignoring")
			return false, nil
		}

		// construct send message api request data
		params := make(map[string]interface{})
		params["message_type"] = ubm.Message.CID.ChatType
		if cid, err := strconv.Atoi(ubm.Message.CID.ChatID); err != nil {
			log.Errorf("[cqhttp-ubm-conv] failed to parse chat id (%v)", err)
			return false, nil
		} else {
			params["user_id"] = cid
			params["group_id"] = cid
			params["discuss_id"] = cid
		}
		params["message"] = message.CQString()
		// params["auto_escape"] = false
		wsRequest := qqbotapi.WebSocketRequest{}
		wsRequest.Action = "send_msg"
		wsRequest.Params = params
		data, _ = json.Marshal(wsRequest)
	}

	// put cqhttp data into packets
	packets := make([]types.Packet, 1)
	packets[0].Head = packet.Head
	packets[0].Head.Format = to
	packets[0].Head.Format.Version = CQHTTPVersion
	packets[0].Body = data

	return true, packets
}
