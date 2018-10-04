package main

type Config struct {
	AdaptorName   string   `toml:"adaptor_name"`
	ChannelUUID   string   `toml:"channel_uuid"`
	CommandPrefix []string `toml:"command_prefix"`
	ResponseMode  uint8    `toml:"response_mode"`
}
