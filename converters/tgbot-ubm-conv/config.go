package main

const (
	UBMAPIVersion         = "1.0"
	TelegramBotAPIVersion = "4.1"
)

type Config struct {
	TelegramAdapters   string `toml:"telegram_adapters"`
	AdapterName        string `toml:"adapter_name"`
	ChannelUUID        string `toml:"channel_uuid"`
	FetchFile          bool   `toml:"fetch_file"`
	APIResponseTimeout string `toml:"api_response_timeout"`
}
