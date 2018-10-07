package main

const (
	UBMAPIVersion = "1.0"
	CQHTTPVersion = "4.5.0"
)

type Config struct {
	QQAdapters         string `toml:"qq_adapters"`
	AdapterName        string `toml:"adapter_name"`
	ChannelUUID        string `toml:"channel_uuid"`
	APIResponseTimeout string `toml:"api_response_timeout"`
	CQHTTPVersion      string `toml:"cqhttp_version"`
}
