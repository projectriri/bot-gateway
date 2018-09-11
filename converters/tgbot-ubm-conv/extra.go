package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/common"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	"github.com/projectriri/bot-gateway/utils"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (p *Plugin) makeRequest(adapter string, endpoint string, values url.Values) []byte {
	endpoint = fmt.Sprintf(APIEndpoint, "00000000:XXXXXXXXXX_XXXXXXXXXXXXXXXXXXXXXXXX", endpoint)
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	req := common.HTTPRequest{
		Method: "POST",
		URL:    endpoint,
		Body:   []byte(values.Encode()),
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
		return resp.Body
	case <-t:
		p.mux.Lock()
		delete(p.pendingRequests, uuid)
		close(ch)
		p.mux.Unlock()
	}
	return nil
}

func (p *Plugin) getFileURL(fileid string, adapter string) string {
	if !p.config.FetchFile {
		return ""
	}
	v := url.Values{}
	v.Add("file_id", fileid)
	resp := p.makeRequest(adapter, "getFile", v)
	if resp == nil {
		return ""
	}
	apiResp := tgbotapi.APIResponse{}
	json.Unmarshal(resp, &apiResp)
	file := tgbotapi.File{}
	json.Unmarshal(apiResp.Result, &file)
	return file.FilePath
}

func (p *Plugin) getMe(adapter string) *ubm_api.User {
	me, ok := p.me[adapter]
	if ok && p.me != nil {
		return me
	}
	v := url.Values{}
	resp := p.makeRequest(adapter, "getMe", v)
	if resp == nil || len(resp) == 0 {
		return nil
	}
	apiResp := tgbotapi.APIResponse{}
	err := json.Unmarshal(resp, &apiResp)
	if err != nil {
		return nil
	}
	user := tgbotapi.User{}
	err = json.Unmarshal(apiResp.Result, &user)
	if err != nil {
		return nil
	}
	u := &ubm_api.User{
		DisplayName: user.FirstName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		UID: ubm_api.UID{
			Messenger: adapter,
			ID:        strconv.Itoa(user.ID),
			Username:  user.UserName,
		},
	}
	p.me[adapter] = u
	return u
}
