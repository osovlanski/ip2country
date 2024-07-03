// cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"ip2country/config"
	"ip2country/internal/api"
	"ip2country/internal/limiter"
	"ip2country/internal/service"
)

func main() {
	// Load configuration.
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Create a rate limiter based on the configuration.
	rateLimiter := limiter.NewRateLimiter(config.RateLimit)

	// Create a service to lookup country information.
	countryService := service.NewCountryService(config.IP2CountryDB)

	// Create and register the HTTP handler for the /v1/find-country endpoint.
	http.HandleFunc("/v1/find-country", api.MakeFindCountryHandler(countryService, rateLimiter))

	// Start the server.
	log.Printf("Starting server on port %s...", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
