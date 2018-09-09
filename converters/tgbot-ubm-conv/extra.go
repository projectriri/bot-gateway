package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/common"
	"github.com/projectriri/bot-gateway/utils"
	"net/http"
	"net/url"
	"time"
)

func (p *Plugin) getFileURL(fileid string, adapter string) string {
	if !p.config.FetchFile {
		return ""
	}
	endpoint := fmt.Sprintf(APIEndpoint, "00000000:XXXXXXXXXX_XXXXXXXXXXXXXXXXXXXXXXXX", "getFile")
	v := url.Values{}
	v.Add("file_id", fileid)
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	req := common.HTTPRequest{
		Method: "POST",
		URL:    endpoint,
		Body:   v.Encode(),
		Header: header,
	}
	uuid := utils.GenerateUUID()
	b, _ := json.Marshal(req)
	packet := types.Packet{
		Head: types.Head{
			From: p.config.AdaptorName,
			To:   adapter,
			UUID: uuid,
			Format: types.Format{
				API:      "telegram-bot-api",
				Version:  "latest",
				Protocol: "http",
				Method:   "apirequest",
			},
		},
		Body: b,
	}
	ch := make(chan types.Packet)
	p.pendingRequests[uuid] = ch
	p.pc.Produce(packet)
	t := time.After(p.timeout)
	select {
	case resp := <-ch:
		apiResp := tgbotapi.APIResponse{}
		json.Unmarshal(resp.Body, &apiResp)
		file := tgbotapi.File{}
		json.Unmarshal(apiResp.Result, &file)
		return file.FilePath
	case <-t:
		p.mux.Lock()
		delete(p.pendingRequests, uuid)
		close(ch)
		p.mux.Unlock()
	}
	return ""
}
