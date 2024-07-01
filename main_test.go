package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"log"
)

func TestFindCountryHandler(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "8082")
	os.Setenv("RATE_LIMIT", "5")
	os.Setenv("IP2COUNTRY_DB", "testdata/ip2country.txt")

	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("could not load config: %v", err)
	}

	log.Printf("Config loaded: %+v", config)

	rateLimiter := NewRateLimiter(config.RateLimit)

	// Create the request
	req, err := http.NewRequest("GET", "/v1/find-country?ip=2.22.233.255", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create the handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		log.Printf("Received request for IP: %s", ip)
		if ip == "" {
			sendErrorResponse(w, "IP is required", http.StatusBadRequest)
			return
		}

		clientIP := r.RemoteAddr
		if !rateLimiter.Allow(clientIP) {
			sendErrorResponse(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		country, city, err := findCountry(ip, config.IP2CountryDB)
		if err != nil {
			sendErrorResponse(w, "IP not found", http.StatusNotFound)
			return
		}

		response := Response{Country: country, City: city}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// Check the response body
	expected := Response{Country: "Israel", City: "Tel-Aviv2"}
	var actual Response
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual, expected)
	}
}


func TestFindCountryHandlerMissingIP(t *testing.T) {
	os.Setenv("PORT", "8082")
	os.Setenv("RATE_LIMIT", "5")
	os.Setenv("IP2COUNTRY_DB", "testdata/ip2country.txt")

	config, _ := LoadConfig()
	rateLimiter := NewRateLimiter(config.RateLimit)

	req, err := http.NewRequest("GET", "/v1/find-country", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			sendErrorResponse(w, "IP is required", http.StatusBadRequest)
			return
		}

		clientIP := r.RemoteAddr
		if !rateLimiter.Allow(clientIP) {
			sendErrorResponse(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		country, city, err := findCountry(ip, config.IP2CountryDB)
		if err != nil {
			sendErrorResponse(w, "IP not found", http.StatusNotFound)
			return
		}

		response := Response{Country: country, City: city}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	handler.ServeHTTP(rr, req)

	expected := ErrorResponse{Error: "IP is required"}
	var actual ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if rr.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	}

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual, expected)
	}
}

func TestFindCountryHandlerNonExistentIP(t *testing.T) {
	os.Setenv("PORT", "8082")
	os.Setenv("RATE_LIMIT", "5")
	os.Setenv("IP2COUNTRY_DB", "testdata/ip2country.txt")

	config, _ := LoadConfig()
	rateLimiter := NewRateLimiter(config.RateLimit)

	req, err := http.NewRequest("GET", "/v1/find-country?ip=1.1.1.1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			sendErrorResponse(w, "IP is required", http.StatusBadRequest)
			return
		}

		clientIP := r.RemoteAddr
		if !rateLimiter.Allow(clientIP) {
			sendErrorResponse(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		country, city, err := findCountry(ip, config.IP2CountryDB)
		if err != nil {
			sendErrorResponse(w, "IP not found", http.StatusNotFound)
			return
		}

		response := Response{Country: country, City: city}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	handler.ServeHTTP(rr, req)

	expected := ErrorResponse{Error: "IP not found"}
	var actual ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if rr.Code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusNotFound)
	}

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual, expected)
	}
}

func TestRateLimitExceeded(t *testing.T) {
	os.Setenv("PORT", "8082")
	os.Setenv("RATE_LIMIT", "1")
	os.Setenv("IP2COUNTRY_DB", "testdata/ip2country.txt")

	config, _ := LoadConfig()
	rateLimiter := NewRateLimiter(config.RateLimit)

	req, err := http.NewRequest("GET", "/v1/find-country?ip=2.22.233.255", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			sendErrorResponse(w, "IP is required", http.StatusBadRequest)
			return
		}

		clientIP := r.RemoteAddr
		if !rateLimiter.Allow(clientIP) {
			sendErrorResponse(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		country, city, err := findCountry(ip, config.IP2CountryDB)
		if err != nil {
			sendErrorResponse(w, "IP not found", http.StatusNotFound)
			return
		}

		response := Response{Country: country, City: city}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// First request should succeed
	handler.ServeHTTP(rr, req)

	// Second request should exceed the rate limit
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected := ErrorResponse{Error: "Rate limit exceeded"}
	var actual ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusTooManyRequests)
	}

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual, expected)
	}
}
