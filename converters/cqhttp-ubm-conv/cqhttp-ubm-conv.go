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
	config         Config
	mux            sync.Mutex
	pc             *router.ProducerChannel
	cc             *router.ConsumerChannel
	me             map[string]*ubm_api.User
	reqChannelPool map[string]chan types.Packet
	timeout        time.Duration
}

var manifest = types.Manifest{
	BasicInfo: types.BasicInfo{
		Name:        "cqhttp-ubm-conv",
		Author:      "Project Riri Staff",
		Version:     "v0.1",
		License:     "MIT",
		URL:         "https://github.com/projectriri/bot-gateway/converters/cqhttp-ubm-conv",
		Description: "Format Converter for CQHTTP API and UBM API.",
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
	var err error
	_, err = toml.DecodeFile(configPath+"/"+filename+".toml", &p.config)
	if err != nil {
		panic(err)
	}
	p.timeout, err = time.ParseDuration(p.config.APIResponseTimeout)
	if err != nil {
		log.Errorf("[cqhttp-ubm-conv] failed to parse api_response_timeout, please check config file")
		panic(err)
	}
	p.me = make(map[string]*ubm_api.User)
	p.reqChannelPool = make(map[string]chan types.Packet)
}

func (p *Plugin) IsConvertible(from types.Format, to types.Format) bool {
	if strings.ToLower(from.API) == "coolq-http-api" && strings.ToLower(to.API) == "ubm-api" {
		if strings.ToLower(from.Method) == "event" && strings.ToLower(to.Method) == "receive" {
			if strings.ToLower(from.Protocol) == "websocket" {
				return true
			}
		}
		if strings.ToLower(from.Method) == "apiresponse" && strings.ToLower(to.Method) == "response" {
			if strings.ToLower(from.Protocol) == "websocket" {
				return true
			}
		}
	}
	if strings.ToLower(from.API) == "ubm-api" && strings.ToLower(to.API) == "coolq-http-api" {
		if strings.ToLower(from.Method) == "send" && strings.ToLower(to.Method) == "apirequest" {
			if strings.ToLower(to.Protocol) == "websocket" {
				return true
			}
		}
	}
	return false
}

func (p *Plugin) Convert(packet types.Packet, to types.Format) (bool, []types.Packet) {
	log.Debugf("[cqhttp-ubm-conv] trying convert pkt %v", packet.Head.UUID)

	from := packet.Head.Format
	if strings.ToLower(from.API) == "coolq-http-api" && strings.ToLower(to.API) == "ubm-api" &&
		utils.CheckIfVersionSatisfy(from.Version, ">=3") && utils.CheckIfVersionSatisfy(UBMAPIVersion, to.Version) {
		if strings.ToLower(from.Method) == "event" && strings.ToLower(to.Method) == "receive" {
			switch strings.ToLower(from.Protocol) {
			case "websocket":
				log.Debugf("[cqhttp-ubm-conv] pkt %v: convertQQEventWSToUbmReceive", packet.Head.UUID)
				return p.convertQQEventWSToUbmReceive(packet, to)
			}
		}
		if strings.ToLower(from.Method) == "apiresponse" && strings.ToLower(to.Method) == "response" {
			switch strings.ToLower(from.Protocol) {
			case "websocket":
				// TODO
			}
		}
	}
	if strings.ToLower(from.API) == "ubm-api" && strings.ToLower(to.API) == "coolq-http-api" &&
		from.Version == "1.0" && utils.CheckIfVersionSatisfy(CQHTTPVersion, to.Version) {
		if strings.ToLower(from.Method) == "send" && strings.ToLower(to.Method) == "apirequest" {
			switch strings.ToLower(to.Protocol) {
			case "websocket":
				log.Debugf("[cqhttp-ubm-conv] pkt %v: convertUbmSendToQQApiRequestWS", packet.Head.UUID)
				return p.convertUbmSendToQQApiRequestWS(packet, to)
			}
		}
	}
	return false, nil
}

func (p *Plugin) Start() {
	log.Infof("[cqhttp-ubm-conv] registering consumer channel %v", p.config.ChannelUUID)
	p.cc = router.RegisterConsumerChannel(p.config.ChannelUUID, []router.RoutingRule{
		{
			From: p.config.QQAdapters,
			To:   p.config.AdapterName,
			Formats: []types.Format{
				{
					API:      "coolq-http-api",
					Version:  ">=3",
					Method:   "apiresponse",
					Protocol: "websocket",
				},
			},
		},
	})
	defer p.cc.Close()
	log.Infof("[cqhttp-ubm-conv] registered consumer channel %v", p.cc.UUID)
	log.Infof("[cqhttp-ubm-conv] registering producer channel %v", p.config.ChannelUUID)
	p.pc = router.RegisterProducerChannel(p.config.ChannelUUID, false)
	defer p.pc.Close()
	for {
		pkt := p.cc.Consume()
		p.mux.Lock()
		ch, ok := p.reqChannelPool[pkt.Head.ReplyToUUID]
		if ok {
			ch <- pkt
			close(ch)
			delete(p.reqChannelPool, pkt.Head.ReplyToUUID)
		}
		p.mux.Unlock()
	}
}

var PluginInstance types.Converter = &Plugin{}
