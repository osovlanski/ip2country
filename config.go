package main

import (
	"log"
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
		log.Println("PORT environment variable not set")
		port = "8080"
	}

	rateLimitStr := os.Getenv("RATE_LIMIT")
	if rateLimitStr == "" {
		log.Println("RATE_LIMIT environment variable not set")
		rateLimitStr = "5"
	}

	rateLimit, err := strconv.Atoi(rateLimitStr)
	if err != nil {
		log.Printf("Invalid RATE_LIMIT value: %v", err)
		return nil, err
	}

	ip2CountryDB := os.Getenv("IP2COUNTRY_DB")
	if ip2CountryDB == "" {
		log.Println("IP2COUNTRY_DB environment variable not set")
		ip2CountryDB = "data/ip2country.txt"
	}

	config := &Config{
		Port:         port,
		RateLimit:    rateLimit,
		IP2CountryDB: ip2CountryDB,
	}

	return config, nil
}
