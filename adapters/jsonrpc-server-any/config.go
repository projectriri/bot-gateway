package main

type Config struct {
	Port            int    `toml:"port"`
	ChannelLifeTime string `toml:"channel_lifetime"`
	GCInterval      string `toml:"garbage_collection_interval"`
}
