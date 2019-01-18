package main

import (
	"io/ioutil"
	gopath "path"
	goplugin "plugin"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	log "github.com/sirupsen/logrus"
	"time"
)

var adapters = make([]types.Adapter, 0)
var converters = make([]types.Converter, 0)
var startTime time.Time

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
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	routerCfg := router.RouterConfig{
		BufferSize: config.BufferSize,
	}

	// init router
	router.Init(routerCfg)

	if config.EnableLittleDaemon {
		InitLittleDaemon()
	}

	startTime = time.Now().Local()

	// load plugins
	ps, err := ioutil.ReadDir(config.PluginDir)
	if err != nil {
		log.Errorf("failed to open types dir %v", config.PluginDir)
	} else {
		for _, p := range ps {
			if gopath.Ext(p.Name()) == ".so" {
				loadPlugin(config.PluginDir + "/" + p.Name())
			}
		}
	}

	if config.EnableLittleDaemon {
		go StartLittleDaemon()
	}

	// start router
	router.Start(converters)

}

func loadPlugin(path string) {
	// get filename without extension
	filenameWithExtension := gopath.Base(path)
	filename := strings.TrimSuffix(filenameWithExtension, gopath.Ext(filenameWithExtension))

	// load .so
	p, err := goplugin.Open(path)
	if err != nil {
		log.Fatalf("failed to load plugin %v: open file error", path)
		return
	}
	// get PluginInstance
	init, err := p.Lookup("Init")
	if err != nil {
		log.Errorf("failed to load plugin %v: Init function not found", path)
		return
	}
	// register and start types
	switch x := init.(type) {
	case func(filename string, configPath string) []types.Adapter:
		adps := x(filename, config.PluginConfDir)
		for _, adp := range adps {
			log.Infof("starting adapter %v", filename)
			go adp.Start()
			adapters = append(adapters, adp)
		}
	case func(filename string, configPath string) []types.Converter:
		covs := x(filename, config.PluginConfDir)
		for _, cov := range covs {
			log.Infof("starting converter %v", filename)
			go cov.Start()
			converters = append(converters, cov)
			adapters = append(adapters, cov)
		}
	default:
		log.Errorf("types %v neither implements an adapter or a converter")
	}
}
