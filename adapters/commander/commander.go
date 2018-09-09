package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/ubm-api"
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
					API:     "ubm-api",
					Version: "1.0",
					Method:  "receive",
				},
			},
		},
	})
	defer cc.Close()
	log.Infof("[commander] registered consumer channel %v", cc.UUID)
	log.Infof("[commander] registering producer channel %v", config.ChannelUUID)
	pc := router.RegisterProducerChannel(config.ChannelUUID, false)
	defer pc.Close()
	log.Infof("[commander] registered producer channel %v", pc.UUID)
	for {
		packet := cc.Consume()
		req := ubm_api.UBM{}
		err := json.Unmarshal(packet.Body, &req)
		if err != nil {
			log.Errorf("[commander] message %v has an incorrect body type %v", packet.Head.UUID, err)
		}
		fmt.Printf("%s\n", string(packet.Body))
		/*ubm := ubm_api.UBM{
			Type: "message",
			Message: &ubm_api.Message{
				CID:      &req.Message.Chat.CID,
				Type:     req.Message.Type,
				ReplyID:  req.Message.ID,
				RichText: req.Message.RichText,
				Location: req.Message.Location,
				Sticker:  req.Message.Sticker,
				Record:   req.Message.Record,
			},
		}
		b, _ := json.Marshal(ubm)
		if req.Message != nil {
			pc.Produce(types.Packet{
				Head: types.Head{
					From:        config.AdaptorName,
					UUID:        utils.GenerateUUID(),
					To:          packet.Head.From,
					ReplyToUUID: packet.Head.UUID,
					Format: types.Format{
						API:      "ubm-api",
						Method:   "send",
						Version:  "1.0",
						Protocol: "",
					},
				},
				Body: b,
			})
		}*/
	}
}

var PluginInstance types.Adapter = &Plugin{}
