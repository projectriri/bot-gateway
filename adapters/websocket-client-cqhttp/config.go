package main

type Config struct {
	AdapterName         string `toml:"adapter_name"`
	ChannelUUID         string `toml:"channel_uuid"`
	CQHTTPAccessToken   string `toml:"cqhttp_access_token"`
	CQHTTPWebSocketAddr string `toml:"cqhttp_websocket_addr"`
	CQHTTPVersion       string `toml:"cqhttp_version"`
}
