package main

import (
	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/types"
	"strings"
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
		Name:    "tgbot-ubm-conv",
		Author:  "Project Riri Staff",
		Version: "v0.1",
		License: "MIT",
		URL:     "https://github.com/projectriri/bot-gateway/converters/tgbot-ubm-conv",
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
	if strings.ToLower(from.API) == "telegram-bot-api" && strings.ToLower(to.API) == "ubm-api" {
		if strings.ToLower(from.Method) == "update" && strings.ToLower(to.Method) == "receive" {
			switch strings.ToLower(from.Protocol) {
			case "http":
				log.Debugf("[tgbot-ubm-conv] pkt %v: convertTgUpdateHttpToUbmReceive", packet.Head.UUID)
				return convertTgUpdateHttpToUbmReceive(packet, to)
			}
		}
		if strings.ToLower(from.Method) == "apiresponse" && strings.ToLower(to.Method) == "response" {
			switch strings.ToLower(from.Protocol) {
			case "http":

			}
		}
	}
	if strings.ToLower(from.API) == "ubm-api" && strings.ToLower(to.API) == "telegram-bot-api" {
		if strings.ToLower(from.Method) == "send" && strings.ToLower(to.Method) == "apirequest" {
			switch strings.ToLower(from.Protocol) {
			case "http":
				log.Debugf("[tgbot-ubm-conv] pkt %v: convertUbmSendToTgApiRequestHttp", packet.Head.UUID)
				return convertUbmSendToTgApiRequestHttp(packet, to)
			}
		}
	}
	return false, nil
}

var PluginInstance types.Converter = &Plugin{}
