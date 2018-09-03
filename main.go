package main

import (
	"time"

	"io/ioutil"
	gopath "path"
	goplugin "plugin"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/plugin"
	"github.com/projectriri/bot-gateway/router"
	log "github.com/sirupsen/logrus"
)

var adaptors = make([]plugin.Adapter, 0)
var converters = make([]plugin.Converter, 0)

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

	// load plugins
	ps, err := ioutil.ReadDir(config.PluginDir)
	if err != nil {
		log.Errorf("failed to open plugin dir %v", config.PluginDir)
	} else {
		for _, p := range ps {
			if gopath.Ext(p.Name()) == ".so" {
				loadPlugin(config.PluginDir + "/" + p.Name())
			}
		}
	}

	// start router
	router.Start(routerCfg)
}

func loadPlugin(path string) {
	// Get filename without extension
	filenameWithExtension := gopath.Base(path)
	filename := strings.TrimSuffix(filenameWithExtension, gopath.Ext(filenameWithExtension))

	// Load .so
	p, err := goplugin.Open(path)
	if err != nil {
		log.Fatalf("failed to load plugin %v: open file error", path)
	}
	// Get PluginInstance
	pi, err := p.Lookup("PluginInstance")
	if err != nil {
		log.Fatalf("failed to load plugin %v: PluginInstance not found", path)
	}
	// Register and Start Plugin
	switch x := pi.(type) {
	case *plugin.Adapter:
		adp := *x
		log.Infof("initializing adapter plugin %v", filename)
		adp.Init(filename)
		log.Infof("starting adapter plugin %v", filename)
		go adp.Start()
		adaptors = append(adaptors, adp)
	case *plugin.Converter:
		cov := *x
		converters = append(converters, cov)
	}
}
