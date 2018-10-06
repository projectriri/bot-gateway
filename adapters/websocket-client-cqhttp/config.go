package main

type Config struct {
	AccessToken         string `toml:"access_token"`
	AdaptorName         string `toml:"adaptor_name"`
	ChannelUUID         string `toml:"channel_uuid"`
	CQHTTPWebSocketAddr string `toml:"cqhttp_websocket_addr"`
}
