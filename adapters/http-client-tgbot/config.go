package main

type Config struct {
	Token       string `toml:"token"`
	AdaptorName string `toml:"adaptor_name"`
	ChannelUUID string `toml:"channel_uuid"`
}
