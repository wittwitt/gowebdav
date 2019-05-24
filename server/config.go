package server

import (
	"github.com/BurntSushi/toml"
)

//LoadConfig config.toml
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

//Config config
type Config struct {
	Listen string
	Prefix string

	Users []*User
}
