package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	RconPassword string `mapstructure:"rcon_password"`
}

type Config struct {
	Server      ServerConfig `mapstructure:"server"`
	GamesMpPath string       `mapstructure:"games_mp_path"`
	RestPort    int          `mapstructure:"rest_port"`
	Environment string       `mapstructure:"environment"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config: %w", err)
	}

	return &config, nil
}
