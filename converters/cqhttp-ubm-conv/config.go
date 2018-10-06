package main

type Config struct {
	QQAdaptors         string `toml:"qq_adaptors"`
	AdaptorName        string `toml:"adaptor_name"`
	ChannelUUID        string `toml:"channel_uuid"`
	APIResponseTimeout string `toml:"api_response_timeout"`
}
