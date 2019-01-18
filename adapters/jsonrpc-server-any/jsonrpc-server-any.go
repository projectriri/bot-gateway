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
	s      *Server
}

var Manifest = types.Manifest{
	BasicInfo: types.BasicInfo{
		Name:        "jsonrpc-server-any",
		Author:      "Project Riri Staff",
		Version:     "v0.1",
		License:     "MIT",
		URL:         "https://github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any",
		Description: "TCP Based JSON RPC Server Adapter.",
	},
	BuildInfo: types.BuildInfo{
		BuildTag:      BuildTag,
		BuildDate:     BuildDate,
		GitCommitSHA1: GitCommitSHA1,
		GitTag:        GitTag,
	},
}

func (p *Plugin) GetManifest() types.Manifest {
	return Manifest
}

func Init(filename string, configPath string) []types.Adapter {
	// load toml config
	var config Config
	_, err := toml.DecodeFile(configPath+"/"+filename+".toml", &config)
	if err != nil {
		panic(err)
	}
	// parse config
	gci, err := time.ParseDuration(config.GCInterval)
	if err != nil {
		log.Error("[jsonrpc-server-any] fail to parse garbage collection interval", err)
		gci = time.Minute * 5
	}
	clt, err := time.ParseDuration(config.ChannelLifeTime)
	if err != nil {
		log.Error("[jsonrpc-server-any] fail to parse channel life time", err)
		clt = time.Hour
	}
	p := Plugin{
		config: config,
	}
	s := new(Server)
	s.init(gci, clt)
	p.s = s
	return []types.Adapter{&p}
}

func (p *Plugin) Start() {
	l, e := net.Listen("tcp", ":"+strconv.Itoa(p.config.Port))
	if e != nil {
		log.Fatalf("[jsonrpc-server-any] listen error: %v", e)
	}
	log.Infof("[jsonrpc-server-any] listening jsonrpc at port %v", p.config.Port)
	go p.s.garbageCollection()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			server := rpc.NewServer()
			b := new(Broker)
			b.init(p.s)
			server.Register(b)
			server.ServeCodec(jsonrpc.NewServerCodec(c))
			b.close()
		}(conn)
	}
}
