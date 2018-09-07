package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

// Telegram constants
const (
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

func (p *Plugin) getUpdates() []byte {

	log.Debug("[longpolling-client-tgbot] pulling updates")

	v := url.Values{}
	if p.updateConfig.Offset != 0 {
		v.Add("offset", strconv.Itoa(p.updateConfig.Offset))
	}
	if p.updateConfig.Limit > 0 {
		v.Add("limit", strconv.Itoa(p.updateConfig.Limit))
	}
	if p.updateConfig.Timeout > 0 {
		v.Add("timeout", strconv.Itoa(p.updateConfig.Timeout))
	}
	method := fmt.Sprintf(APIEndpoint, p.config.Token, "getUpdates")

	resp, err := p.client.PostForm(method, v)
	defer resp.Body.Close()
	if err != nil {
		return nil
	}

	// read response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	log.Debugf("[longpolling-client-tgbot] %v", string(data))

	// update offset id
	var apiResp tgbotapi.APIResponse
	err = json.Unmarshal(data, &apiResp)
	if err != nil {
		return nil
	}
	var updates []tgbotapi.Update
	json.Unmarshal(apiResp.Result, &updates)
	for _, update := range updates {
		if update.UpdateID >= p.updateConfig.Offset {
			p.updateConfig.Offset = update.UpdateID + 1
		}
	}

	if len(updates) > 0 {
		return data
	} else {
		return nil
	}
}
