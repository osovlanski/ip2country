// cmd/server/main.go
package main

import (
	"net/http"
	"ip2country/config"
	"ip2country/internal/api"
	"ip2country/internal/limiter"
	"ip2country/internal/service"

	"github.com/sirupsen/logrus"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("could not load config: %v", err)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})

	rateLimiter := limiter.NewRateLimiter(config.RateLimit)
	countryService := service.NewCountryService(config.IP2CountryDB)

	http.HandleFunc("/v1/find-country", api.MakeFindCountryHandler(countryService, rateLimiter))

	logrus.Infof("Starting server on port %s...", config.Port)
	logrus.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
