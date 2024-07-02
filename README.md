# IP to Country Service

This project implements a simple IP to country service. The service receives an HTTP GET request with an IP address and returns the corresponding country and city. The service also includes a rate-limiting mechanism and can be extended to support multiple IP to country databases.

## Features

- Clear and easy-to-read code
- Configuration via environment variables using Viper
- HTTP GET endpoint to retrieve country and city by IP address
- JSON responses for both success and error cases
- Rate-limiting mechanism without using external libraries
- Easily extendable to support multiple IP to country databases
- Dockerized for easy deployment

## Getting Started

### Prerequisites

- Go 1.18 or later
- Docker

### Setup

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/ip2country.git
    cd ip2country
    ```

2. Create a `data` directory and add your IP to country database file:

    ```sh
    mkdir data
    echo "2.22.233.255,City,Country" > data/ip2country.txt
    ```

3. Set the required environment variables in a `.env` file:

    ```sh
    echo "PORT=8080" >> .env
    echo "RATE_LIMIT=5" >> .env
    echo "IP2COUNTRY_DB=data/ip2country.txt" >> .env
    ```

### Running Locally

1. Install Go dependencies:

    ```sh
    go mod download
    ```

2. Build and run the application:

    ```sh
    go build -o main ./cmd/server
    ./main
    ```

3. Send a test request:

    ```sh
    curl "http://localhost:8080/v1/find-country?ip=2.22.233.255"
    ```

### Running with Docker

1. Build the Docker image:

    ```sh
    docker build -t ip2country:latest .
    ```

2. Run the Docker containers using Docker Compose:

    ```sh
    docker-compose up --build
    ```

3. Send a test request:

    ```sh
    curl "http://localhost:8080/v1/find-country?ip=2.22.233.255"
    ```

### Testing

1. Run tests:

    ```sh
    go test ./...
    ```

### Pushing to Docker Hub

1. Log in to Docker Hub:

    ```sh
    docker login
    ```

2. Tag the Docker image:

    ```sh
    docker tag ip2country:latest yourdockerhubusername/ip2country:latest
    ```

3. Push the Docker image to Docker Hub:

    ```sh
    docker push yourdockerhubusername/ip2country:latest
    ```

## Environment Variables

- `PORT`: The port on which the server will run.
- `RATE_LIMIT`: The maximum number of requests per second allowed.
- `IP2COUNTRY_DB`: The path to the IP to country database file.

## API Endpoints

- `GET /v1/find-country?ip=<ip>`: Returns the country and city for the given IP address.

## Example Response

### Success

```json
{
    "country": "Country",
    "city": "City"
}
