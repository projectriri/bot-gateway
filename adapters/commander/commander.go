package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/ubm-api"
	log "github.com/sirupsen/logrus"
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
		Name:    "commander",
		Author:  "Project Riri Staff",
		Version: "v0.1",
		License: "MIT",
		URL:     "https://github.com/projectriri/bot-gateway/adapters/commander",
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
	log.Infof("[commander] registering consumer channel %v", config.ChannelUUID)
	cc := router.RegisterConsumerChannel(config.ChannelUUID, []router.RoutingRule{
		{
			From: ".*",
			To:   ".*",
			Formats: []types.Format{
				{
					API:     "UBM-API",
					Version: "1.0",
					Method:  "Receive",
				},
			},
		},
	})
	log.Infof("[commander] registered consumer channel %v", cc.UUID)
	log.Infof("[commander] registering producer channel %v", config.ChannelUUID)
	pc := router.RegisterProducerChannel(config.ChannelUUID, false)
	log.Infof("[commander] registered producer channel %v", pc.UUID)
	for {
		packet := cc.Consume()
		req, ok := packet.Body.(*ubm_api.UBM)
		if !ok {
			log.Errorf("[commander] message %v has an incorrect body type", packet.Head.UUID)
		}
		fmt.Printf("%+v\n", req)
	}
}

var PluginInstance types.Adapter = &Plugin{}
