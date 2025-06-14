package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Address string        `yaml:"address"`
	Plugins PluginsConfig `yaml:"plugins"`
}

type PluginsConfig struct {
	Prefix     string     `yaml:"prefix"`
	PrefixFile string     `yaml:"prefixFile"`
	CORS       CORSConfig `yaml:"cors"`
}

type CORSConfig struct {
	AllowOrigins []string `yaml:"allowOrigins"`
}

func LoadConfig(path string) (*Config, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	if err := yaml.Unmarshal(buf, conf); err != nil {
		return nil, err
	}

	if conf.Address == "" {
		conf.Address = ":8081"
	}

	if conf.Plugins.PrefixFile != "" {
		p, err := os.ReadFile(conf.Plugins.PrefixFile)
		if err != nil {
			return nil, err
		}
		conf.Plugins.Prefix = string(p)
	}

	return conf, nil
}
