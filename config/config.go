package config

import (
	"github.com/BurntSushi/toml"
	"github.com/spf13/pflag"
)

type Config struct {
	BotToken string  `toml:"bot_token"`
	Admin    []int64 `toml:"admin"`
	GroupId  int64   `toml:"group_id"`
	Proxy    Proxy   `toml:"proxy"`
	Debug    bool    `toml:"debug"`
}

type Proxy struct {
	URL string `toml:"url"`
}

var (
	c *Config
)

func LoadConfig() error {
	configPath := pflag.StringP("config", "c", "config.toml", "Path to the configuration file")
	pflag.Parse()
	c = new(Config)
	_, err := toml.DecodeFile(*configPath, c)
	return err
}

func GetConfig() *Config {
	return c
}
