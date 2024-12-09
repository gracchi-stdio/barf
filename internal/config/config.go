package config

import "github.com/spf13/viper"

type Server struct {
	Port string
}

type App struct {
	ENV string
}

type Config struct {
	Server Server
}

func Load() (*Config, error) {

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("app.env", "development")

	viper.SetConfigName("config")

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
