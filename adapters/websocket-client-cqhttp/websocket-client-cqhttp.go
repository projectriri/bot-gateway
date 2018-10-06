package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/websocket"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
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
		Name:    "websocket-client-cqhttp",
		Author:  "Project Riri Staff",
		Version: "v0.1",
		License: "MIT",
		URL:     "https://github.com/projectriri/bot-gateway/adapters/websocket-client-cqhttp",
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
}

func (p *Plugin) Start() {
	log.Infof("[websocket-client-cqhttp] registering consumer channel %v", p.config.ChannelUUID)
	cc := router.RegisterConsumerChannel(p.config.ChannelUUID, []router.RoutingRule{
		{
			From: ".*",
			To:   p.config.AdaptorName,
			Formats: []types.Format{
				{
					API:      "coolq-http-api",
					Version:  p.config.CQHTTPVersion,
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
	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Token %s", p.config.CQHTTPAccessToken))
	// Dial /api/ ws
	p.apiClient, _, err = websocket.DefaultDialer.Dial(
		p.config.CQHTTPWebSocketAddr+"/api/",
		header,
	)
	if err != nil {
		log.Errorf("[websocket-client-cqhttp] failed to dial cqhttp api websocket (%v)", err)
	} else {
		log.Infof("[websocket-client-cqhttp] dial cqhttp api websocket success")
	}
	defer p.apiClient.Close()
	// Dial /event/ ws
	p.eventClient, _, err = websocket.DefaultDialer.Dial(
		p.config.CQHTTPWebSocketAddr+"/event/",
		header,
	)
	if err != nil {
		log.Errorf("[websocket-client-cqhttp] failed to dial cqhttp event websocket (%v)", err)
	} else {
		log.Infof("[websocket-client-cqhttp] dial cqhttp event websocket success")
	}
	defer p.eventClient.Close()

	// Start main event update loop
	go func() {
		for {
			_, msg, err := p.eventClient.ReadMessage()
			if err != nil {
				log.Errorf("[websocket-client-cqhttp] failed to read event (%v)", err)
				continue
			}
			log.Debugf("[websocket-client-cqhttp] receiving event %s", string(msg))
			pc.Produce(types.Packet{
				Head: types.Head{
					From: p.config.AdaptorName,
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

	// Start main api request loop
	go func() {
		for {
			// send api request
			apiRequestPkt := cc.Consume()
			err := p.apiClient.WriteMessage(websocket.TextMessage, apiRequestPkt.Body)
			if err != nil {
				log.Errorf("[websocket-client-cqhttp] failed to send apirequest (%v)", err)
				continue
			}
			// get api response
			_, msg, err := p.apiClient.ReadMessage()
			if err != nil {
				log.Errorf("[websocket-client-cqhttp] failed to read apiresponse (%v)", err)
				continue
			}
			log.Debugf("[websocket-client-cqhttp] receiving apiresponse %s", string(msg))
			pc.Produce(types.Packet{
				Head: types.Head{
					From:        p.config.AdaptorName,
					To:          apiRequestPkt.Head.From,
					UUID:        utils.GenerateUUID(),
					ReplyToUUID: apiRequestPkt.Head.UUID,
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

	// lock the main thread
	<-make(chan bool)
}

var PluginInstance types.Adapter = &Plugin{}
