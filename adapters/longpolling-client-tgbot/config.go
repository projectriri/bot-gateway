package main

type Config struct {
	Token       string `toml:"token"`
	Limit       int    `toml:"limit"`
	Timeout     int    `toml:"timeout"`
	AdaptorName string `toml:"adaptor_name"`
	ChannelUUID string `toml:"channel_uuid"`
}
