package main

import (
	"time"
	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/router"
)

func main() {

	// Load Global Config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		panic(err)
	}

	// Parse Router Config
	gci, _ := time.ParseDuration(config.GCInterval)
	clt, err := time.ParseDuration(config.ChannelLifeTime)
	if err != nil {
		// log.Error(err)
		clt = time.Hour
	}
	routerCfg := router.RouterConfig {
		BufferSize: config.BufferSize,
		ChannelLifeTime: clt,
		GCInterval: gci,
	}

	// Start Router
	router.Start(routerCfg)
}
