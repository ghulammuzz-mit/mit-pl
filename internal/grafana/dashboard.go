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

// Dashboard represents a Grafana dashboard summary.
type Dashboard struct {
	UID     string `json:"uid"`
	Title   string `json:"title"`
	Folder  string `json:"folderTitle"`
	Tags    []string `json:"tags"`
	Updated int64  `json:"updated"`
}

// ListDashboardsResponse represents the response from listing dashboards.
type ListDashboardsResponse struct {
	Title       string `json:"title"`
	UID         string `json:"uid"`
	FolderTitle string `json:"folderTitle"`
	FolderUID   string `json:"folderUid"`
	FolderURL   string `json:"folderUrl"`
	Type        string `json:"type"`
	Tags        []string `json:"tags"`
	IsStarred   bool   `json:"isStarred"`
	URL         string `json:"url"`
}

// SearchDashboardsParams represents parameters for searching dashboards.
type SearchDashboardsParams struct {
	Query string // Search query string
	Type  string // Dashboard type (e.g., "dash-db")
	Tags  []string // Filter by tags
	Limit int    // Maximum number of results
}

// ListDashboards returns all available dashboards.
func (c *Client) ListDashboards(ctx context.Context) ([]Dashboard, error) {
	return c.SearchDashboards(ctx, &SearchDashboardsParams{})
}

// SearchDashboards searches for dashboards based on the provided parameters.
func (c *Client) SearchDashboards(ctx context.Context, params *SearchDashboardsParams) ([]Dashboard, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client configuration: %w", err)
	}

	// Build query URL
	baseURL, err := url.Parse(c.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid Grafana URL: %w", err)
	}

	apiURL := baseURL.ResolveReference(&url.URL{
		Path: "/api/search",
	})

	// Build query parameters
	values := url.Values{}
	if params != nil {
		if params.Query != "" {
			values.Add("query", params.Query)
		}
		if params.Type != "" {
			values.Add("type", params.Type)
		} else {
			values.Add("type", "dash-db") // Default to dashboards
		}
		for _, tag := range params.Tags {
			values.Add("tag", tag)
		}
		if params.Limit > 0 {
			values.Add("limit", fmt.Sprintf("%d", params.Limit))
		}
	}

	apiURL.RawQuery = values.Encode()

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	slog.Info("Searching Grafana dashboards", "url", apiURL.String())

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
	var results []ListDashboardsResponse
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to Dashboard type
	dashboards := make([]Dashboard, 0, len(results))
	for _, r := range results {
		dashboards = append(dashboards, Dashboard{
			UID:     r.UID,
			Title:   r.Title,
			Folder:  r.FolderTitle,
			Tags:    r.Tags,
			Updated: 0, // Not included in search response
		})
	}

	slog.Info("Retrieved dashboards", "count", len(dashboards))

	return dashboards, nil
}

// GetDashboard retrieves a specific dashboard by UID.
func (c *Client) GetDashboard(ctx context.Context, uid string) (map[string]interface{}, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client configuration: %w", err)
	}

	if uid == "" {
		return nil, fmt.Errorf("dashboard UID is required")
	}

	// Build query URL
	baseURL, err := url.Parse(c.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid Grafana URL: %w", err)
	}

	apiURL := baseURL.ResolveReference(&url.URL{
		Path: "/api/dashboards/uid/" + uid,
	})

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	slog.Info("Getting Grafana dashboard", "uid", uid, "url", apiURL.String())

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

	// Parse response - return raw map for flexibility
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}

