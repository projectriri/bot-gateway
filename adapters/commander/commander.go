package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/cmd"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	"github.com/projectriri/bot-gateway/utils"
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
		log.Warnf("%s\n", string(packet.Body))
		if req.Type == "message" && req.Message != nil {
			if req.Message.Type != "rich_text" || req.Message.RichText == nil {
				continue
			}
			slices := make([][]ubm_api.RichTextElement, 0)

			// TODO: Deal with command prefix

			// TODO: Process this Message
			for _, elem := range *req.Message.RichText {
				if elem.Type == "text" {

				}
			}

			c := cmd.Command{
				Cmd:  slices[0],
				// TODO: CmdStr
				Args: slices[1:],
				// TODO: ArgsTxt
				// TODO: ArgsStr
			}

			b, _ := json.Marshal(c)
			pc.Produce(types.Packet{
				Head: types.Head{
					From: config.AdaptorName,
					UUID: utils.GenerateUUID(),
					Format: types.Format{
						API:      "cmd",
						Method:   "cmd",
						Version:  "1.0",
						Protocol: "",
					},
				},
				Body: b,
			})
		}

	}
}

var PluginInstance types.Adapter = &Plugin{}
