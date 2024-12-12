package config

import (
	"github.com/spf13/viper"
	"strings"
)

type Server struct {
	Port string
	Env  string
}

type Database struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

type Config struct {
	Server Server
	DB     Database
}

func Load() (*Config, error) {
	// set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.env", "development")

	// database defaults
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.user", "")
	viper.SetDefault("db.password", "")
	viper.SetDefault("db.name", "")
	viper.SetDefault("db.port", "5432")

	// config file settings
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// environment variable settings
	//viper.SetEnvPrefix("BARF")                             // all env vars will be prefixed with BARF
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // replace . with _
	viper.AutomaticEnv()

	// read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// override with environment variables
	if envPort := viper.GetString("PORT"); envPort != "" {
		cfg.Server.Port = envPort
	}
	if dbHost := viper.GetString("DB_HOST"); dbHost != "" {
		cfg.DB.Host = dbHost
	}
	if dbUser := viper.GetString("DB_USER"); dbUser != "" {
		cfg.DB.User = dbUser
	}
	if dbPassword := viper.GetString("DB_PASSWORD"); dbPassword != "" {
		cfg.DB.Password = dbPassword
	}
	if dbName := viper.GetString("DB_NAME"); dbName != "" {
		cfg.DB.Name = dbName
	}
	if dbPort := viper.GetString("DB_PORT"); dbPort != "" {
		cfg.DB.Port = dbPort
	}

	return &cfg, nil
}
