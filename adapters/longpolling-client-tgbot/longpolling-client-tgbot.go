package main

import (
	"github.com/BurntSushi/toml"
	"github.com/go-telegram-bot-api/telegram-bot-api"
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
	client       http.Client
	updateConfig tgbotapi.UpdateConfig
	config       Config
}

var manifest = types.Manifest{
	BasicInfo: types.BasicInfo{
		Name:        "longpolling-client-tgbot",
		Author:      "Project Riri Staff",
		Version:     "v0.1",
		License:     "MIT",
		URL:         "https://github.com/projectriri/bot-gateway/adapters/longpolling-client-tgbot",
		Description: "Long Polling Client Adapter for Telegram Bot API.",
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
	p.updateConfig.Limit = p.config.Limit
	p.updateConfig.Timeout = p.config.Timeout
}

func (p *Plugin) Start() {
	log.Infof("[longpolling-client-tgbot] registering producer channel %v", p.config.ChannelUUID)
	pc := router.RegisterProducerChannel(p.config.ChannelUUID, false)
	defer pc.Close()
	log.Infof("[longpolling-client-tgbot] registered producer channel %v", pc.UUID)
	log.Info("[longpolling-client-tgbot] start polling from Telegram-Bot-API via LongPolling")
	var data []byte
	for {
		data = p.getUpdates()
		if data != nil {
			log.Debug("[longpolling-client-tgbot] producing packet")
			pc.Produce(types.Packet{
				Head: types.Head{
					From: p.config.AdapterName,
					UUID: utils.GenerateUUID(),
					Format: types.Format{
						API:      "telegram-bot-api",
						Version:  APIVersion,
						Method:   "update",
						Protocol: "http",
					},
				},
				Body: data,
			})
		}
	}
}

var PluginInstance types.Adapter = &Plugin{}
