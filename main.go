package main

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/router"
	log "github.com/sirupsen/logrus"
)

func main() {

	// load global config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		panic(err)
	}

	// init logger
	lv, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		lv = log.InfoLevel
	}
	log.SetLevel(lv)

	// parse router config
	gci, err := time.ParseDuration(config.GCInterval)
	if err != nil {
		log.Error("fail to parse garbage collection interval", err)
		gci = time.Minute * 5
	}
	clt, err := time.ParseDuration(config.ChannelLifeTime)
	if err != nil {
		log.Error("fail to parse channel life time", err)
		clt = time.Hour
	}
	routerCfg := router.RouterConfig{
		BufferSize:      config.BufferSize,
		ChannelLifeTime: clt,
		GCInterval:      gci,
	}

	// start router
	router.Start(routerCfg)
}
