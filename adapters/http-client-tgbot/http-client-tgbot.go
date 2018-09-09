package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/common"
	"github.com/projectriri/bot-gateway/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var (
	BuildTag      string
	BuildDate     string
	GitCommitSHA1 string
	GitTag        string
)

type Plugin struct {
	client http.Client
	config Config
}

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
	_, err := toml.DecodeFile(configPath+"/"+filename+".toml", &p.config)
	if err != nil {
		panic(err)
	}
}

func (p *Plugin) Start() {
	log.Infof("[http-client-tgbot] registering consumer channel %v", p.config.ChannelUUID)
	cc := router.RegisterConsumerChannel(p.config.ChannelUUID, []router.RoutingRule{
		{
			From: ".*",
			To:   p.config.AdaptorName,
			Formats: []types.Format{
				{
					API:      "telegram-bot-api",
					Version:  "latest",
					Method:   "apirequest",
					Protocol: "http",
				},
			},
		},
	})
	defer cc.Close()
	log.Infof("[http-client-tgbot] registered consumer channel %v", cc.UUID)
	log.Infof("[http-client-tgbot] registering producer channel %v", p.config.ChannelUUID)
	pc := router.RegisterProducerChannel(p.config.ChannelUUID, false)
	defer pc.Close()
	log.Infof("[http-client-tgbot] registered producer channel %v", pc.UUID)
	for {
		packet := cc.Consume()
		r := common.HTTPRequest{}
		err := json.Unmarshal(packet.Body, &r)
		if err != nil {
			log.Errorf("[http-client-tgbot] message %v has an incorrect body type %v", packet.Head.UUID, err)
			continue
		}
		req, err := http.NewRequest(r.Method, r.URL, strings.NewReader(r.Body))
		if err != nil {
			log.Errorf("[http-client-tgbot] message %v has an incorrect body type %v", packet.Head.UUID, err)
			continue
		}
		req.Header = r.Header
		data, err := p.makeRequest(req)
		if err != nil {
			continue
		}
		pc.Produce(types.Packet{
			Head: types.Head{
				From:        p.config.AdaptorName,
				To:          packet.Head.From,
				UUID:        utils.GenerateUUID(),
				ReplyToUUID: packet.Head.UUID,
				Format: types.Format{
					API:      "telegram-bot-api",
					Version:  "latest",
					Method:   "apiresponse",
					Protocol: "http",
				},
			},
			Body: data,
		})
	}
}

var PluginInstance types.Adapter = &Plugin{}
