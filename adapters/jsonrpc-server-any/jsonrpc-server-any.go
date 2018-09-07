package main

import (
	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/types"
	log "github.com/sirupsen/logrus"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"time"
)

var (
	BuildTag      string
	BuildDate     string
	GitCommitSHA1 string
	GitTag        string
)

type Plugin struct {
	config Config
	b      Broker
}

var manifest = types.Manifest{
	BasicInfo: types.BasicInfo{
		Name:    "jsonrpc-server-any",
		Author:  "Project Riri Staff",
		Version: "v0.1",
		License: "MIT",
		URL:     "https://github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any",
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
	_, err := toml.DecodeFile(configPath+"/"+filename+".toml", &p.config)
	if err != nil {
		panic(err)
	}
	// parse config
	gci, err := time.ParseDuration(p.config.GCInterval)
	if err != nil {
		log.Error("fail to parse garbage collection interval", err)
		gci = time.Minute * 5
	}
	clt, err := time.ParseDuration(p.config.ChannelLifeTime)
	if err != nil {
		log.Error("fail to parse channel life time", err)
		clt = time.Hour
	}
	b := new(Broker)
	b.init(gci, clt)
}

func (p *Plugin) Start() {
	server := rpc.NewServer()
	server.Register(p.b)
	l, e := net.Listen("tcp", ":"+strconv.Itoa(p.config.Port))
	if e != nil {
		log.Fatalf("[jsonrpc-server-any] listen error: %v", e)
	}
	log.Infof("[jsonrpc-server-any] listening jsonrpc at port %v", p.config.Port)
	go p.b.garbageCollection()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

var PluginInstance types.Adapter = &Plugin{}
