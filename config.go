package main

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string
	RateLimit    int
	IP2CountryDB string
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	rateLimitStr := os.Getenv("RATE_LIMIT")
	if rateLimitStr == "" {
		rateLimitStr = "5"
	}

	rateLimit, err := strconv.Atoi(rateLimitStr)
	if err != nil {
		return nil, err
	}

	ip2countryDB := os.Getenv("IP2COUNTRY_DB")
	if ip2countryDB == "" {
		ip2countryDB = "data/ip2country.txt"
	}

	return &Config{
		Port:         port,
		RateLimit:    rateLimit,
		IP2CountryDB: ip2countryDB,
	}, nil
}
