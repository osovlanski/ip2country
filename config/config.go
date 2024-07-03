// config/config.go
package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds the configuration values for the application.
type Config struct {
	Port         string // The port the server will listen on.
	RateLimit    int    // The maximum number of requests per second allowed.
	IP2CountryDB string // The path to the IP to country database file.
}

// LoadConfig loads configuration from environment variables file.
// If values are not found, it uses default values.
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
