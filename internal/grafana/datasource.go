package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

// DataSource represents a Grafana data source.
type DataSource struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	URL   string `json:"url"`
	JSONData struct {
		HTTPMethod string `json:"httpMethod"`
	} `json:"jsonData"`
	IsDefault bool `json:"isDefault"`
}

// ListDataSourcesResponse represents the response from listing data sources.
type ListDataSourcesResponse struct {
	DataSources []DataSource `json:"dataSources"`
}

// ListDataSources returns all available data sources.
func (c *Client) ListDataSources(ctx context.Context) ([]DataSource, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client configuration: %w", err)
	}

	// Build query URL
	baseURL, err := url.Parse(c.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid Grafana URL: %w", err)
	}

	apiURL := baseURL.ResolveReference(&url.URL{
		Path: "/api/datasources",
	})

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	slog.Info("Listing Grafana data sources", "url", apiURL.String(), "api_key_prefix", c.APIKey[:10]+"...", "full_url", apiURL.String())

	// Execute request
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	// Parse response
	var result []DataSource
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	slog.Info("Retrieved data sources", "count", len(result))

	return result, nil
}

// GetDataSource retrieves a specific data source by UID.
func (c *Client) GetDataSource(ctx context.Context, uid string) (*DataSource, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client configuration: %w", err)
	}

	if uid == "" {
		return nil, fmt.Errorf("data source UID is required")
	}

	// Build query URL
	baseURL, err := url.Parse(c.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid Grafana URL: %w", err)
	}

	apiURL := baseURL.ResolveReference(&url.URL{
		Path: "/api/datasources/uid/" + uid,
	})

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	slog.Info("Getting Grafana data source", "uid", uid, "url", apiURL.String())

	// Execute request
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	// Parse response
	var result DataSource
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// FindPrometheusDataSources finds all Prometheus-type data sources.
func (c *Client) FindPrometheusDataSources(ctx context.Context) ([]DataSource, error) {
	sources, err := c.ListDataSources(ctx)
	if err != nil {
		return nil, err
	}

	var promSources []DataSource
	for _, source := range sources {
		if source.Type == "prometheus" {
			promSources = append(promSources, source)
		}
	}

	slog.Info("Found Prometheus data sources", "count", len(promSources))

	return promSources, nil
}

// GetDefaultDataSource returns the default data source.
func (c *Client) GetDefaultDataSource(ctx context.Context) (*DataSource, error) {
	sources, err := c.ListDataSources(ctx)
	if err != nil {
		return nil, err
	}

	for _, source := range sources {
		if source.IsDefault {
			slog.Info("Found default data source", "uid", source.UID, "name", source.Name)
			return &source, nil
		}
	}

	return nil, fmt.Errorf("no default data source found")
}
