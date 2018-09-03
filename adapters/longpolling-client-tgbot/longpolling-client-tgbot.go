package main

import (
	"github.com/projectriri/bot-gateway/plugin"
	"fmt"
)

var (
	BuildTag      string
	BuildDate     string
	GitCommitSHA1 string
	GitTag        string
)

type Plugin struct{}

var manifest = plugin.Manifest{
	BasicInfo: plugin.BasicInfo{
		Name:    "longpolling-client-tgbot",
		Author:  "Project Riri Stuff",
		Version: "v0.1",
		License: "MIT",
		URL:     "https://github.com/projectriri/bot-gateway/adapters/longpolling-client-tgbot",
	},
	BuildInfo: plugin.BuildInfo{
		BuildTag:      BuildTag,
		BuildDate:     BuildDate,
		GitCommitSHA1: GitCommitSHA1,
		GitTag:        GitTag,
	},
}

func (p *Plugin) GetManifest() plugin.Manifest {
	return manifest
}

func (p *Plugin) Init(filename string) {
	fmt.Printf("Hey let's init %v\n", filename)
}

func (p *Plugin) Start() {
	fmt.Printf("Sudāto!!!!!\n")
}

var PluginInstance plugin.Adapter = &Plugin{}
