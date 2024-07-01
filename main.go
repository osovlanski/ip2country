package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	rateLimiter := NewRateLimiter(config.RateLimit)

	http.HandleFunc("/v1/find-country", func(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Starting server on port %s...", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

func findCountry(ip, dbPath string) (string, string, error) {
	file, err := os.Open(dbPath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue
		}
		if parts[0] == ip {
			return parts[2], parts[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", err
	}

	return "", "", errors.New("IP not found")
}
