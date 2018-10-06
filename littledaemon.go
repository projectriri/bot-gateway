package main

import (
	"encoding/json"
	"fmt"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/cmd"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	"github.com/projectriri/bot-gateway/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type D struct {
	Day    string `yaml:"day"`
	Hour   string `yaml:"hour"`
	Minute string `yaml:"minute"`
	Second string `yaml:"second"`
}

type T struct {
	Template string `yaml:"template"`
	Duration D      `yaml:"duration"`
}

var t T

func InitLittleDaemon() {
	b, err := ioutil.ReadFile("locale.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(b, &t)
	if err != nil {
		log.Fatalf("error parsing yaml file: %v", err)
	}
}

func StartLittleDaemon() {
	log.Infof("[littledaemon] registering consumer channel %s", config.LittleDaemonChannelUUID)
	cc := router.RegisterConsumerChannel(config.LittleDaemonChannelUUID, []router.RoutingRule{
		{
			From: ".*",
			To:   ".*",
			Formats: []types.Format{
				{
					API:     "cmd",
					Version: "1.0",
					Method:  "cmd",
				},
			},
		},
	})
	defer cc.Close()
	log.Infof("[littledaemon] registered consumer channel %v", cc.UUID)
	log.Infof("[littledaemon] registering producer channel %v", config.LittleDaemonChannelUUID)
	pc := router.RegisterProducerChannel(config.LittleDaemonChannelUUID, false)
	defer pc.Close()
	log.Infof("[littledaemon] registered producer channel %v", pc.UUID)
	for {
		packet := cc.Consume()
		req := cmd.Command{}
		err := json.Unmarshal(packet.Body, &req)
		if err != nil {
			log.Errorf("[littledaemon] message %v has an incorrect body type %v", packet.Head.UUID, err)
		}

		if req.CmdStr == "ping" || (len(req.Cmd) > 0 && req.Cmd[0].Text == "ping") ||
			req.CmdStr == "status" || (len(req.Cmd) > 0 && req.Cmd[0].Text == "status") {
			now := time.Now().Local()
			i, o := router.GetIOCount()
			str := fmt.Sprintf(
				t.Template,
				now.Format("2006-01-02 15:04:05"),
				router.GetChannelCount(),
				router.GetCachedPacketCount(),
				router.GetChannelCacheLimit(),
				getPluginCount(),
				formatDuration(now.Sub(startTime)),
				i, o,
			)
			resp := ubm_api.UBM{
				Type: "message",
				Message: &ubm_api.Message{
					Type: "rich_text",
					RichText: &ubm_api.RichText{
						{
							Type: "text",
							Text: str,
						},
					},
					CID:     &req.Message.Chat.CID,
					ReplyID: req.Message.ID,
				},
			}
			b, _ := json.Marshal(resp)
			pc.Produce(types.Packet{
				Head: types.Head{
					From: config.LittleDaemonName,
					To:   packet.Head.From,
					UUID: utils.GenerateUUID(),
					Format: types.Format{
						API:      "ubm-api",
						Method:   "send",
						Version:  "1.0",
						Protocol: "",
					},
				},
				Body: b,
			})
		}
	}
}

func formatDuration(duration time.Duration) string {
	d := duration / (24 * time.Hour)
	h1 := duration / time.Hour
	h := h1 - d*24
	m1 := duration / time.Minute
	m := m1 - h1*60
	s := duration/time.Second - m1*60
	msg := ""
	flag := false
	if d > 0 {
		msg += fmt.Sprintf("%d %s ", d, t.Duration.Day)
		flag = true
	}
	if h > 0 || flag {
		msg += fmt.Sprintf("%d %s ", h, t.Duration.Hour)
		flag = true
	}
	if m > 0 || flag {
		msg += fmt.Sprintf("%d %s ", m, t.Duration.Minute)
		flag = true
	}
	msg += fmt.Sprintf("%d %s", s, t.Duration.Second)
	return msg
}

func getPluginCount() int {
	return len(adapters)
}
