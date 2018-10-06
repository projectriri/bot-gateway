package main

type Config struct {
	Token       string `toml:"token"`
	AdapterName string `toml:"adapter_name"`
	ChannelUUID string `toml:"channel_uuid"`
}
