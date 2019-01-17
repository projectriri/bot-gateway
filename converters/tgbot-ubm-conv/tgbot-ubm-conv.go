package main

import (
	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	"github.com/projectriri/bot-gateway/utils"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

var (
	BuildTag      string
	BuildDate     string
	GitCommitSHA1 string
	GitTag        string
)

type Plugin struct {
	pendingRequests map[string]chan types.Packet
	mux             sync.Mutex
	config          Config
	pc              *router.ProducerChannel
	cc              *router.ConsumerChannel
	timeout         time.Duration
	me              map[string]*ubm_api.User
}

var manifest = types.Manifest{
	BasicInfo: types.BasicInfo{
		Name:        "tgbot-ubm-conv",
		Author:      "Project Riri Staff",
		Version:     "v0.1",
		License:     "MIT",
		URL:         "https://github.com/projectriri/bot-gateway/converters/tgbot-ubm-conv",
		Description: "Format Converter for Telegram Bot API and UBM API.",
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

// Telegram constants
const (
	APIVersion = "4.1"
)

func (p *Plugin) Init(filename string, configPath string) {
	// load toml config
	_, err := toml.DecodeFile(configPath+"/"+filename+".toml", &p.config)
	if err != nil {
		panic(err)
	}
	p.pendingRequests = make(map[string]chan types.Packet)
	p.timeout, err = time.ParseDuration(p.config.APIResponseTimeout)
	if err != nil {
		log.Error("[tgbot-ubm-conv] fail to parse timeout", err)
		p.timeout = time.Minute * 5
	}
	p.me = make(map[string]*ubm_api.User)
}

func (p *Plugin) IsConvertible(from types.Format, to types.Format) bool {
	if strings.ToLower(from.API) == "telegram-bot-api" && strings.ToLower(to.API) == "ubm-api" {
		if strings.ToLower(from.Method) == "update" && strings.ToLower(to.Method) == "receive" {
			if strings.ToLower(from.Protocol) == "http" {
				return true
			}
		}
		if strings.ToLower(from.Method) == "apiresponse" && strings.ToLower(to.Method) == "response" {
			if strings.ToLower(from.Protocol) == "http" {
				return true
			}
		}
	}
	if strings.ToLower(from.API) == "ubm-api" && strings.ToLower(to.API) == "telegram-bot-api" {
		if strings.ToLower(from.Method) == "send" && strings.ToLower(to.Method) == "apirequest" {
			if strings.ToLower(to.Protocol) == "http" {
				return true
			}
		}
	}
	return false
}

func (p *Plugin) Convert(packet types.Packet, to types.Format) (bool, []types.Packet) {
	log.Debugf("[tgbot-ubm-conv] try convert pkt %v", packet.Head.UUID)

	from := packet.Head.Format
	if strings.ToLower(from.API) == "telegram-bot-api" && strings.ToLower(to.API) == "ubm-api" &&
		utils.CheckIfVersionSatisfy(from.Version, ">=3") && utils.CheckIfVersionSatisfy(UBMAPIVersion, to.Version) {
		if strings.ToLower(from.Method) == "update" && strings.ToLower(to.Method) == "receive" {
			switch strings.ToLower(from.Protocol) {
			case "http":
				log.Debugf("[tgbot-ubm-conv] pkt %v: convertTgUpdateHttpToUbmReceive", packet.Head.UUID)
				return p.convertTgUpdateHttpToUbmReceive(packet, to)
			}
		}
		if strings.ToLower(from.Method) == "apiresponse" && strings.ToLower(to.Method) == "response" {
			switch strings.ToLower(from.Protocol) {
			case "http":
				// TODO
			}
		}
	}
	if strings.ToLower(from.API) == "ubm-api" && strings.ToLower(to.API) == "telegram-bot-api" &&
		from.Version == "1.0" && utils.CheckIfVersionSatisfy(TelegramBotAPIVersion, to.Version) {
		if strings.ToLower(from.Method) == "send" && strings.ToLower(to.Method) == "apirequest" {
			switch strings.ToLower(to.Protocol) {
			case "http":
				log.Debugf("[tgbot-ubm-conv] pkt %v: convertUbmSendToTgApiRequestHttp", packet.Head.UUID)
				return p.convertUbmSendToTgApiRequestHttp(packet, to)
			}
		}
	}
	return false, nil
}

func (p *Plugin) Start() {
	log.Infof("[tgbot-ubm-conv] registering consumer channel %v", p.config.ChannelUUID)
	p.cc = router.RegisterConsumerChannel(p.config.ChannelUUID, []router.RoutingRule{
		{
			From: p.config.TelegramAdapters,
			To:   p.config.AdapterName,
			Formats: []types.Format{
				{
					API:      "telegram-bot-api",
					Version:  ">=3",
					Method:   "apiresponse",
					Protocol: "http",
				},
			},
		},
	})
	defer p.cc.Close()
	log.Infof("[tgbot-ubm-conv] registered consumer channel %v", p.cc.UUID)
	log.Infof("[tgbot-ubm-conv] registering producer channel %v", p.config.ChannelUUID)
	p.pc = router.RegisterProducerChannel(p.config.ChannelUUID, false)
	defer p.pc.Close()
	for pkt := range p.cc.Buffer {
		p.mux.Lock()
		ch, ok := p.pendingRequests[pkt.Head.ReplyToUUID]
		if ok {
			ch <- pkt
			close(ch)
			delete(p.pendingRequests, pkt.Head.ReplyToUUID)
		}
		p.mux.Unlock()
	}
}

var PluginInstance types.Converter = &Plugin{}
