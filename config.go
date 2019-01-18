package main

var config = globalConfig{}

type globalConfig struct {
	BufferSize              uint   `toml:"buffer_size"`
	LogLevel                string `toml:"log_level"`
	PluginDir               string `toml:"plugin_dir"`
	PluginConfDir           string `toml:"plugin_conf_dir"`
}
