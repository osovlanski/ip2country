// test/api/handler_test.go
package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"ip2country/config"
	"ip2country/internal/api"
	"ip2country/internal/limiter"
	"ip2country/internal/service"
	"ip2country/pkg"
)

// Helper function to create the handler
func createHandler() (http.Handler, error) {
	// Load configuration
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	rateLimiter := limiter.NewRateLimiter(config.RateLimit)
	countryService := service.NewCountryService(config.IP2CountryDB)
	handler := api.MakeFindCountryHandler(countryService, rateLimiter)

	return handler, nil
}

// Test main function to set environment variables
func TestMain(m *testing.M) {
	// Set environment variables for testing
	os.Setenv("PORT", "8080")
	os.Setenv("RATE_LIMIT", "5")
	os.Setenv("IP2COUNTRY_DB", "/app/testdata/ip2country.txt")

	code := m.Run()
	os.Exit(code)
}

func TestFindCountryHandler(t *testing.T) {
	handler, err := createHandler()
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("GET", "/v1/find-country?ip=2.22.233.255", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	// Check the response body
	expected := pkg.Response{Country: "Country", City: "City"}
	var actual pkg.Response
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestFindCountryHandlerMissingIP(t *testing.T) {
	handler, err := createHandler()
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	req, err := http.NewRequest("GET", "/v1/find-country", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected := pkg.ErrorResponse{Error: "IP is required"}
	var actual pkg.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if rr.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
	}

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestFindCountryHandlerNonExistentIP(t *testing.T) {
	handler, err := createHandler()
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	req, err := http.NewRequest("GET", "/v1/find-country?ip=1.1.1.1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected := pkg.ErrorResponse{Error: "IP not found"}
	var actual pkg.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if rr.Code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusNotFound)
	}

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestRateLimitExceeded(t *testing.T) {
	// Set the environment variable for rate limit to 1 for this test
	os.Setenv("RATE_LIMIT", "1")

	handler, err := createHandler()
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	req, err := http.NewRequest("GET", "/v1/find-country?ip=2.22.233.255", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// First request should succeed
	handler.ServeHTTP(rr, req)

	// Second request should exceed the rate limit
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected := pkg.ErrorResponse{Error: "Rate limit exceeded"}
	var actual pkg.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusTooManyRequests)
	}

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}

	// Reset the rate limit environment variable
	os.Setenv("RATE_LIMIT", "5")
}
