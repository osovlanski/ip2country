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
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	rateLimiter := limiter.NewRateLimiter(config.RateLimit)
	countryService := service.NewCountryService(config.IP2CountryDB)

	http.HandleFunc("/v1/find-country", api.MakeFindCountryHandler(countryService, rateLimiter))

	log.Printf("Starting server on port %s...", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
