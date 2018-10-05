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

var adaptors = make([]types.Adapter, 0)
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
		log.Fatalf("failed to load types %v: open file error", path)
		return
	}
	// get PluginInstance
	pi, err := p.Lookup("PluginInstance")
	if err != nil {
		log.Errorf("failed to load types %v: PluginInstance not found", path)
		return
	}
	// register and start types
	switch x := pi.(type) {
	case *types.Adapter:
		adp := *x
		log.Infof("initializing adapter types %v", filename)
		adp.Init(filename, config.PluginConfDir)
		log.Infof("starting adapter types %v", filename)
		go adp.Start()
		adaptors = append(adaptors, adp)
	case *types.Converter:
		cov := *x
		converters = append(converters, cov)
		log.Infof("initializing adapter types %v", filename)
		cov.Init(filename, config.PluginConfDir)
		log.Infof("starting adapter types %v", filename)
		go cov.Start()
		adaptors = append(adaptors, cov)
	default:
		log.Errorf("types %v neither implements an adapter or a converter")
	}
}
