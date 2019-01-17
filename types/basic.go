package types

type BasicInfo struct {
	Name        string
	Version     string
	Author      string
	License     string
	URL         string
	Description string
}

type BuildInfo struct {
	BuildTag      string
	BuildDate     string
	GitCommitSHA1 string
	GitTag        string
}

type Manifest struct {
	BasicInfo BasicInfo
	BuildInfo BuildInfo
}

type BasePlugin interface {
	GetManifest() Manifest
	Init(filename string, configPath string)
}
