package main

type Config struct {
	TelegramAdaptors string `toml:"telegram_adaptors"`
	AdaptorName      string `toml:"adaptor_name"`
	ChannelUUID      string `toml:"channel_uuid"`
	FetchFile        bool   `toml:"fetch_file"`
	FetchFileTimeout string `toml:"fetch_file_timeout"`
}
