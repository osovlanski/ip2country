// config/config.go
package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Port         string
	RateLimit    int
	IP2CountryDB string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Warnf("Error reading config file: %v", err)
	}

	config := &Config{
		Port:         viper.GetString("PORT"),
		RateLimit:    viper.GetInt("RATE_LIMIT"),
		IP2CountryDB: viper.GetString("IP2COUNTRY_DB"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	if config.RateLimit == 0 {
		config.RateLimit = 5
	}

	if config.IP2CountryDB == "" {
		config.IP2CountryDB = "data/ip2country.txt"
	}

	return config, nil
}
