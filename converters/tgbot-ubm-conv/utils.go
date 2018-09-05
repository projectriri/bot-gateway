package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
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

func newMessageRequest(endpoint string, params url.Values) *http.Request {
	endpoint = fmt.Sprintf(APIEndpoint, "00000000:XXXXXXXXXX_XXXXXXXXXXXXXXXXXXXXXXXX", endpoint)
	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}
