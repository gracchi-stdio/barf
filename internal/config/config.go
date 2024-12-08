package config

import "github.com/spf13/viper"

type Server struct {
	Port string
}
type Config struct {
	Server Server
}

func Load() (*Config, error) {

	viper.SetDefault("server.port", "8080")

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
