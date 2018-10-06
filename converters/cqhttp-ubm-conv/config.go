package main

type Config struct {
	QQAdapters         string `toml:"qq_adapters"`
	AdapterName        string `toml:"adapter_name"`
	ChannelUUID        string `toml:"channel_uuid"`
	APIResponseTimeout string `toml:"api_response_timeout"`
}
