package main

var config = globalConfig{}

type globalConfig struct {
	BufferSize              uint   `toml:"buffer_size"`
	LogLevel                string `toml:"log_level"`
	PluginDir               string `toml:"plugin_dir"`
	PluginConfDir           string `toml:"plugin_conf_dir"`
	EnableLittleDaemon      bool   `toml:"enable_little_daemon"`
	LittleDaemonName        string `toml:"little_daemon_name"`
	LittleDaemonChannelUUID string `toml:"little_daemon_channel_uuid"`
}
