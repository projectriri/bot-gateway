package main

import (
	"bytes"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/projectriri/bot-gateway/types/common"
	"mime/multipart"
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

type PhotoConfig struct {
	Type  string `json:"type"`
	Media string `json:"media"`
}

func getTelegramChatType(chat *tgbotapi.Chat) string {
	if chat.IsSuperGroup() {
		return "supergroup"
	} else if chat.IsGroup() {
		return "group"
	} else if chat.IsChannel() {
		return "channel"
	} else if chat.IsPrivate() {
		return "private"
	} else {
		return ""
	}
}

func newMessageRequest(endpoint string, params map[string]string) common.HTTPRequest {
	endpoint = fmt.Sprintf(APIEndpoint, "00000000:XXXXXXXXXX_XXXXXXXXXXXXXXXXXXXXXXXX", endpoint)
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	req := common.HTTPRequest{
		Method: "POST",
		URL:    endpoint,
		Body:   values.Encode(),
		Header: header,
	}
	return req
}

func newFileRequest(endpoint string, params map[string]string, files map[string][]byte) common.HTTPRequest {
	if len(files) == 0 {
		return newMessageRequest(endpoint, params)
	}
	endpoint = fmt.Sprintf(APIEndpoint, "00000000:XXXXXXXXXX_XXXXXXXXXXXXXXXXXXXXXXXX", endpoint)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	ct := "multipart/form-data; boundary=" + writer.Boundary()
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	for k, v := range files {
		part, err := writer.CreateFormFile(k, k)
		if err != nil {
			fmt.Println(err)
		}
		part.Write(v)
	}
	writer.Close()
	header := http.Header{}
	header.Add("Content-Type", ct)
	req := common.HTTPRequest{
		Method: "POST",
		URL:    endpoint,
		Body:   body.String(),
		Header: header,
	}
	return req
}

func plainToMarkdown(from string) (to string) {
	from = strings.Replace(from, `\`, `\\`, -1)
	from = strings.Replace(from, `_`, `\_`, -1)
	from = strings.Replace(from, `*`, `\*`, -1)
	from = strings.Replace(from, "`", "\\`", -1)
	return from
}
