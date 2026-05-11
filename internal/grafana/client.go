package grafana

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// Client is a Grafana API client.
type Client struct {
	HTTP   *http.Client
	URL    string
	APIKey string
}

// New creates a new Grafana client from environment variables.
// It reads GRAFANA_URL and GRAFANA_API_KEY from the environment.
func New() *Client {
	url := os.Getenv("GRAFANA_URL")
	apiKey := os.Getenv("GRAFANA_API_KEY")

	// Debug logging
	fmt.Fprintf(os.Stderr, "DEBUG: GRAFANA_URL=%s\n", url)
	fmt.Fprintf(os.Stderr, "DEBUG: GRAFANA_API_KEY length=%d\n", len(apiKey))
	if len(apiKey) > 10 {
		fmt.Fprintf(os.Stderr, "DEBUG: GRAFANA_API_KEY prefix=%s...\n", apiKey[:10])
	}

	return &Client{
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
		URL:    url,
		APIKey: apiKey,
	}
}

// Validate checks if the client has the required configuration.
func (c *Client) Validate() error {
	if c.URL == "" {
		return &ConfigError{Field: "GRAFANA_URL", Message: "Grafana URL is required"}
	}
	if c.APIKey == "" {
		return &ConfigError{Field: "GRAFANA_API_KEY", Message: "Grafana API key is required"}
	}
	return nil
}

// ConfigError represents a configuration error.
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Field + ": " + e.Message
}

// APIError represents an error from the Grafana API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}
