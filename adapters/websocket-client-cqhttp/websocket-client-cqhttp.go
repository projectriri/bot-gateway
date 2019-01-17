package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/catsworld/qq-bot-api"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

var (
	BuildTag      string
	BuildDate     string
	GitCommitSHA1 string
	GitTag        string
)

type Plugin struct {
	apiClient   *websocket.Conn
	eventClient *websocket.Conn
	config      Config
}

var manifest = types.Manifest{
	BasicInfo: types.BasicInfo{
		Name:        "websocket-client-cqhttp",
		Author:      "Project Riri Staff",
		Version:     "v0.1",
		License:     "MIT",
		URL:         "https://github.com/projectriri/bot-gateway/adapters/websocket-client-cqhttp",
		Description: "WebSocket Client Adapter for CoolQ HTTP API.",
	},
	BuildInfo: types.BuildInfo{
		BuildTag:      BuildTag,
		BuildDate:     BuildDate,
		GitCommitSHA1: GitCommitSHA1,
		GitTag:        GitTag,
	},
}

func (p *Plugin) GetManifest() types.Manifest {
	return manifest
}

func (p *Plugin) Init(filename string, configPath string) {
	// load toml config
	_, err := toml.DecodeFile(configPath+"/"+filename+".toml", &p.config)
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Errorf("[cqhttp-ubm-conv] failed to parse api_response_timeout, please check config file")
		panic(err)
	}
}

func (p *Plugin) Start() {
	log.Infof("[websocket-client-cqhttp] registering consumer channel %v", p.config.ChannelUUID)
	cc := router.RegisterConsumerChannel(p.config.ChannelUUID, []router.RoutingRule{
		{
			From: ".*",
			To:   p.config.AdapterName,
			Formats: []types.Format{
				{
					API:      "coolq-http-api",
					Version:  ">=3",
					Method:   "apirequest",
					Protocol: "websocket",
				},
			},
		},
	})
	defer cc.Close()
	log.Infof("[websocket-client-cqhttp] registered consumer channel %v", cc.UUID)

	log.Infof("[websocket-client-cqhttp] registering producer channel %v", p.config.ChannelUUID)
	pc := router.RegisterProducerChannel(p.config.ChannelUUID, false)
	defer pc.Close()
	log.Infof("[websocket-client-cqhttp] registered producer channel %v", pc.UUID)

	log.Infof("[websocket-client-cqhttp] dialing cqhttp-websocket server")
	var err error
	// Dial /api/ ws
	apiConfig, err := websocket.NewConfig(p.config.CQHTTPWebSocketAddr+"/api/", "http://localhost/")
	if err != nil {
		log.Fatalf("[websocket-client-cqhttp] invalid websocket address %v", err)
	}
	apiConfig.Header.Add("Authorization", fmt.Sprintf("Token %s", p.config.CQHTTPAccessToken))
	p.apiClient, err = websocket.DialConfig(apiConfig)
	if err != nil {
		log.Errorf("[websocket-client-cqhttp] failed to dial cqhttp api websocket (%v)", err)
	} else {
		log.Infof("[websocket-client-cqhttp] dial cqhttp api websocket success")
	}
	defer p.apiClient.Close()
	// Dial /event/ ws
	eventConfig, err := websocket.NewConfig(p.config.CQHTTPWebSocketAddr+"/event/", "http://localhost/")
	if err != nil {
		log.Fatalf("[websocket-client-cqhttp] invalid websocket address %v", err)
	}
	eventConfig.Header.Add("Authorization", fmt.Sprintf("Token %s", p.config.CQHTTPAccessToken))
	p.eventClient, err = websocket.DialConfig(eventConfig)
	if err != nil {
		log.Errorf("[websocket-client-cqhttp] failed to dial cqhttp event websocket (%v)", err)
	} else {
		log.Infof("[websocket-client-cqhttp] dial cqhttp event websocket success")
	}
	defer p.eventClient.Close()

	// Start main event update loop
	go func() {
		for {
			msg := json.RawMessage{}
			if err := websocket.JSON.Receive(p.eventClient, &msg); err != nil {
				log.Errorf("[websocket-client-cqhttp] failed to read event (%v)", err)
				continue
			}
			log.Debugf("[websocket-client-cqhttp] receiving event %s", string(msg))
			pc.Produce(types.Packet{
				Head: types.Head{
					From: p.config.AdapterName,
					To:   "",
					UUID: utils.GenerateUUID(),
					Format: types.Format{
						API:      "coolq-http-api",
						Version:  p.config.CQHTTPVersion,
						Method:   "event",
						Protocol: "websocket",
					},
				},
				Body: msg,
			})
		}
	}()

	go func() {
		for {
			// get api response
			msg := json.RawMessage{}
			if err := websocket.JSON.Receive(p.apiClient, &msg); err != nil {
				log.Errorf("[websocket-client-cqhttp] failed to read apiresponse (%v)", err)
				continue
			}
			apiResp := qqbotapi.APIResponse{}
			err := json.Unmarshal(msg, &apiResp)
			if err != nil {
				log.Errorf("[websocket-client-cqhttp] failed to parse apiresponse (%v)", err)
				continue
			}
			e, ok := apiResp.Echo.(string)
			if !ok {
				log.Errorf("[websocket-client-cqhttp] bad echo %v", apiResp)
				continue
			}

			if len(e) < 36 || !utils.ValidateUUID(e[:36]) {
				log.Errorf("[websocket-client-cqhttp] bad echo %v", apiResp)
				continue
			}

			log.Debugf("[websocket-client-cqhttp] receiving apiresponse %s", string(msg))
			pc.Produce(types.Packet{
				Head: types.Head{
					From:        p.config.AdapterName,
					To:          e[36:],
					UUID:        utils.GenerateUUID(),
					ReplyToUUID: e[:36],
					Format: types.Format{
						API:      "coolq-http-api",
						Version:  p.config.CQHTTPVersion,
						Method:   "apiresponse",
						Protocol: "websocket",
					},
				},
				Body: msg,
			})
		}
	}()

	// Start main api request loop
	go func() {
		for {
			// send api request
			apiRequestPkt := cc.Consume()
			apiReq := qqbotapi.WebSocketRequest{}
			err := json.Unmarshal(apiRequestPkt.Body, &apiReq)
			if err != nil {
				log.Errorf("[websocket-client-cqhttp] failed to parse apirequest (%v)", err)
				continue
			}

			apiReq.Echo = apiRequestPkt.Head.UUID + apiRequestPkt.Head.From
			err = websocket.JSON.Send(p.apiClient, apiReq)
			if err != nil {
				log.Errorf("[websocket-client-cqhttp] failed to send apirequest (%v)", err)
				continue
			}
		}
	}()

	// lock the main thread
	<-make(chan bool)
}

var PluginInstance types.Adapter = &Plugin{}
