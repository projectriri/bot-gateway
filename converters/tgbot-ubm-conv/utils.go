package main

import "github.com/go-telegram-bot-api/telegram-bot-api"

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
