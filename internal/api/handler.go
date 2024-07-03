// internal/api/handler.go
package api

import (
	"encoding/json"
	"net/http"
	"ip2country/internal/limiter"
	"ip2country/internal/service"
	"ip2country/pkg"
)

// MakeFindCountryHandler creates an HTTP handler for the required endpoint (e.g. /v1/find-country).
// It accepts a CountryService and a RateLimiter as parameters.
func MakeFindCountryHandler(countryService *service.CountryService, rateLimiter *limiter.RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			pkg.SendErrorResponse(w, "IP is required", http.StatusBadRequest)
			return
		}

		clientIP := r.RemoteAddr
		if !rateLimiter.Allow(clientIP) {
			pkg.SendErrorResponse(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		country, city, err := countryService.FindCountry(ip)
		if err != nil {
			pkg.SendErrorResponse(w, "IP not found", http.StatusNotFound)
			return
		}

		response := pkg.Response{Country: country, City: city}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
