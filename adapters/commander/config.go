package main

var config = Config{}

type Config struct {
	AdaptorName string `toml:"adaptor_name"`
	ChannelUUID string `toml:"channel_uuid"`
}
