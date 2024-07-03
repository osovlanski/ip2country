// internal/service/country_service.go
package service

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

type CountryService struct {
	dbPath string
}

// NewCountryService creates a new CountryService with the given database path.
func NewCountryService(dbPath string) *CountryService {
	return &CountryService{dbPath: dbPath}
}

// FindCountry finds the country and city for a given IP address.
// It returns the country and city as strings, or an error if the IP is not found.
func (s *CountryService) FindCountry(ip string) (string, string, error) {
	file, err := os.Open(s.dbPath)
	if err != nil {
		log.Printf("Failed to open database file: %v", err)
		return "", "", err
	}
	defer file.Close()

	// Scan the file line by line to find the matching IP address.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue
		}
		if parts[0] == ip {
			log.Printf("Found country: %s, city: %s for IP: %s", parts[2], parts[1], ip)
			return parts[2], parts[1], nil
		}
	}

	// Check for errors during scanning.
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading database file: %v", err)
		return "", "", err
	}

	log.Printf("IP not found in database: %s", ip)
	return "", "", errors.New("IP not found")
}
