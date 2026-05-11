package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

// QueryRequest represents a query request to Grafana.
type QueryRequest struct {
	DataSource string            `json:"dataSource"` // UID of the data source
	Query      string            `json:"query"`      // PromQL query
	StartTime  time.Time         `json:"startTime"`
	EndTime    time.Time         `json:"endTime"`
	Step       time.Duration     `json:"step"` // Resolution step
	Labels     map[string]string `json:"labels,omitempty"`
}

// QueryResponse represents the response from a Grafana query.
type QueryResponse struct {
	Data struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"` // [timestamp, value]
		} `json:"result"`
	} `json:"data"`
	Status string `json:"status"`
}

// Query executes a PromQL query against Grafana.
func (c *Client) Query(ctx context.Context, req *QueryRequest) (*QueryResponse, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client configuration: %w", err)
	}

	if req.DataSource == "" {
		return nil, fmt.Errorf("data source UID is required")
	}
	if req.Query == "" {
		return nil, fmt.Errorf("query is required")
	}

	// Build query URL for Grafana's datasource proxy API
	baseURL, err := url.Parse(c.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid Grafana URL: %w", err)
	}

	queryURL := baseURL.ResolveReference(&url.URL{
		Path: fmt.Sprintf("/api/datasources/proxy/uid/%s/api/v1/query_range", req.DataSource),
	})

	// Build query parameters
	values := url.Values{}
	values.Add("query", req.Query)
	values.Add("start", fmt.Sprintf("%d", req.StartTime.Unix()))
	values.Add("end", fmt.Sprintf("%d", req.EndTime.Unix()))
	values.Add("step", fmt.Sprintf("%.0f", req.Step.Seconds()))

	queryURL.RawQuery = values.Encode()

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "GET", queryURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	slog.Info("Executing Grafana query",
		"url", queryURL.String(),
		"query", req.Query,
		"start", req.StartTime,
		"end", req.EndTime,
		"step", req.Step,
	)

	// Execute request
	resp, err := c.HTTP.Do(httpReq)
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
		return nil, &QueryError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	// Parse response
	var result QueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// QueryInstant executes an instant query (returns single value).
func (c *Client) QueryInstant(ctx context.Context, dataSource, query string, time time.Time) (*QueryResponse, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client configuration: %w", err)
	}

	if dataSource == "" {
		return nil, fmt.Errorf("data source UID is required")
	}
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	// Build query URL for Grafana's datasource proxy API
	baseURL, err := url.Parse(c.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid Grafana URL: %w", err)
	}

	queryURL := baseURL.ResolveReference(&url.URL{
		Path: fmt.Sprintf("/api/datasources/proxy/uid/%s/api/v1/query", dataSource),
	})

	// Build query parameters
	values := url.Values{}
	values.Add("query", query)
	values.Add("time", fmt.Sprintf("%d", time.Unix()))

	queryURL.RawQuery = values.Encode()

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "GET", queryURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	slog.Info("Executing instant Grafana query",
		"url", queryURL.String(),
		"query", query,
		"time", time,
	)

	// Execute request
	resp, err := c.HTTP.Do(httpReq)
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
		return nil, &QueryError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	// Parse response
	var result QueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// QueryError represents an error from a query operation.
type QueryError struct {
	StatusCode int
	Message    string
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("query failed with status %d: %s", e.StatusCode, e.Message)
}

// DirectQuery executes a query using Grafana's direct query API.
// This uses /api/datasources/proxy/uid/{uid}/api/v1/query_range
func (c *Client) DirectQuery(ctx context.Context, dataSource, query string, start, end time.Time, step time.Duration) (*QueryResponse, error) {
	return c.Query(ctx, &QueryRequest{
		DataSource: dataSource,
		Query:      query,
		StartTime:  start,
		EndTime:    end,
		Step:       step,
	})
}

// QueryRaw executes a raw query against Grafana's query API.
// This is useful for debugging or complex queries.
func (c *Client) QueryRaw(ctx context.Context, dataSource string, queryExpr string, start, end time.Time, step time.Duration) (json.RawMessage, error) {
	req := &QueryRequest{
		DataSource: dataSource,
		Query:      queryExpr,
		StartTime:  start,
		EndTime:    end,
		Step:       step,
	}

	baseURL, err := url.Parse(c.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid Grafana URL: %w", err)
	}

	queryURL := baseURL.ResolveReference(&url.URL{
		Path: fmt.Sprintf("/api/datasources/proxy/uid/%s/api/v1/query_range", dataSource),
	})

	values := url.Values{}
	values.Add("query", req.Query)
	values.Add("start", fmt.Sprintf("%d", req.StartTime.Unix()))
	values.Add("end", fmt.Sprintf("%d", req.EndTime.Unix()))
	values.Add("step", fmt.Sprintf("%.0f", req.Step.Seconds()))

	queryURL.RawQuery = values.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, "GET", queryURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &QueryError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	return body, nil
}
