package main

import (
	"github.com/BurntSushi/toml"
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

type Plugin struct{}

var manifest = types.Manifest{
	BasicInfo: types.BasicInfo{
		Name:    "http-client-tgbot",
		Author:  "Project Riri Staff",
		Version: "v0.1",
		License: "MIT",
		URL:     "https://github.com/projectriri/bot-gateway/adapters/http-client-tgbot",
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
	_, err := toml.DecodeFile(configPath+"/"+filename+".toml", &config)
	if err != nil {
		panic(err)
	}
}

func (p *Plugin) Start() {
	log.Infof("[http-client-tgbot] registering consumer channel %v", config.ChannelUUID)
	cc := router.RegisterConsumerChannel(config.ChannelUUID, []router.RoutingRule{
		{
			From: ".*",
			To:   config.AdaptorName,
			Formats: []types.Format{
				{
					API:      "Telegram-Bot-API",
					Version:  "latest",
					Method:   "APIRequest",
					Protocol: "HTTP",
				},
			},
		},
	})
	log.Infof("[http-client-tgbot] registered consumer channel %v", cc.UUID)
	log.Infof("[http-client-tgbot] registering producer channel %v", config.ChannelUUID)
	pc := router.RegisterProducerChannel(config.ChannelUUID, false)
	log.Infof("[http-client-tgbot] registered producer channel %v", pc.UUID)
	for {
		packet := cc.Consume()
		req, ok := packet.Body.(*http.Request)
		if !ok {
			log.Errorf("[http-client-tgbot] message %v has an incorrect body type", packet.Head.UUID)
		}
		data, err := makeRequest(req)
		if err != nil {
			continue
		}
		pc.Produce(types.Packet{
			Head: types.Head{
				From:        config.AdaptorName,
				UUID:        utils.GenerateUUID(),
				ReplyToUUID: packet.Head.UUID,
				Format: types.Format{
					API:      "Telegram-Bot-API",
					Version:  "latest",
					Method:   "APIResponse",
					Protocol: "HTTP",
				},
			},
			Body: data,
		})
	}
}

var PluginInstance types.Adapter = &Plugin{}
