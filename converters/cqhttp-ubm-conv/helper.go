package main

import (
	"encoding/json"
	"github.com/catsworld/qq-bot-api"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	"github.com/projectriri/bot-gateway/utils"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func (p *Plugin) makeRequest(to string, action string, params map[string]interface{}) []byte {
	// construct api request packet
	request := make(map[string]interface{})
	request["action"] = action
	request["params"] = params
	body, _ := json.Marshal(request)
	pkt := types.Packet{
		Head: types.Head{
			UUID: utils.GenerateUUID(),
			From: p.config.AdaptorName,
			To:   to,
			Format: types.Format{
				API:      "coolq-http-api",
				Version:  "latest",
				Method:   "apirequest",
				Protocol: "websocket",
			},
		},
		Body: body,
	}

	// make channel for api response
	p.reqChannelPool[pkt.Head.UUID] = make(chan types.Packet)

	// push packet into router
	p.pc.Produce(pkt)

	// wait for api response or timeout
	timeout := time.After(p.timeout)
	select {
	case res := <-p.reqChannelPool[pkt.Head.UUID]:
		// unmarshal api response
		response := qqbotapi.APIResponse{}
		if err := json.Unmarshal(res.Body, &response); err != nil {
			log.Errorf("[cqhttp-ubm-conv] fail to parse api response json")
			return nil
		}

		// check status
		if response.Status != "ok" {
			log.Errorf(
				"[cqhttp-ubm-conv] api response status: %s, retcode: %d",
				response.Status,
				response.RetCode,
			)
			return nil
		}
		return response.Data
	case <-timeout:
		p.mux.Lock()
		close(p.reqChannelPool[pkt.Head.UUID])
		delete(p.reqChannelPool, pkt.Head.UUID)
		p.mux.Unlock()
		log.Errorf("[cqhttp-ubm-conv] waiting for api response"+
			" timeout (req pkt UUID: %s)", pkt.Head.UUID)
		return nil
	}
	return nil
}

func (p *Plugin) getMe(adapter string) *ubm_api.User {
	// check cache
	if u, ok := p.me[adapter]; ok && u != nil {
		return u
	}

	// cache miss, make request
	if res := p.makeRequest(adapter, "get_login_info", nil); res != nil {
		qqUser := qqbotapi.User{}
		if err := json.Unmarshal(res, &qqUser); err != nil {
			log.Errorf("[cqhttp-ubm-conv] fail to parse get_login_info json")
			return nil
		}
		user := ubm_api.User{
			DisplayName: qqUser.NickName,
			FirstName:   qqUser.NickName,
			UID: ubm_api.UID{
				Messenger: adapter,
				ID:        strconv.FormatInt(qqUser.ID, 10),
			},
		}
		// save to cache
		p.me[adapter] = &user
		return &user
	}
	return nil
}
