package tools

import (
	"fmt"
	"strings"
	"time"
)

// PromQLBuilder helps build PromQL queries from parameters.
type PromQLBuilder struct {
	metric     string
	labels     map[string]string
	aggregation string
	percentile float64
	timeWindow time.Duration
}

// NewPromQLBuilder creates a new PromQL builder.
func NewPromQLBuilder(metric string) *PromQLBuilder {
	return &PromQLBuilder{
		metric:     metric,
		labels:     make(map[string]string),
		aggregation: "avg",
	}
}

// WithLabel adds a label filter.
func (b *PromQLBuilder) WithLabel(key, value string) *PromQLBuilder {
	b.labels[key] = value
	return b
}

// WithLabels adds multiple label filters.
func (b *PromQLBuilder) WithLabels(labels map[string]string) *PromQLBuilder {
	for k, v := range labels {
		b.labels[k] = v
	}
	return b
}

// WithAggregation sets the aggregation type.
func (b *PromQLBuilder) WithAggregation(aggregation string) *PromQLBuilder {
	b.aggregation = aggregation
	return b
}

// WithPercentile sets the percentile value.
func (b *PromQLBuilder) WithPercentile(p float64) *PromQLBuilder {
	b.percentile = p
	b.aggregation = "percentile"
	return b
}

// WithTimeWindow sets the time window for rate/irate queries.
func (b *PromQLBuilder) WithTimeWindow(d time.Duration) *PromQLBuilder {
	b.timeWindow = d
	return b
}

// Build returns the constructed PromQL query.
func (b *PromQLBuilder) Build() string {
	// Start with metric name and labels
	query := b.metric

	if len(b.labels) > 0 {
		labels := make([]string, 0, len(b.labels))
		for k, v := range b.labels {
			labels = append(labels, fmt.Sprintf(`%s="%s"`, k, v))
		}
		query += "{" + strings.Join(labels, ",") + "}"
	}

	// Wrap with aggregation
	switch b.aggregation {
	case "percentile":
		if b.percentile > 0 {
			query = fmt.Sprintf(`quantile(%.2f, %s)`, b.percentile/100, query)
		}
	case "avg":
		query = fmt.Sprintf(`avg(%s)`, query)
	case "sum":
		query = fmt.Sprintf(`sum(%s)`, query)
	case "max":
		query = fmt.Sprintf(`max(%s)`, query)
	case "min":
		query = fmt.Sprintf(`min(%s)`, query)
	case "count":
		query = fmt.Sprintf(`count(%s)`, query)
	}

	return query
}

// BuildRate builds a rate query with the specified time window.
func (b *PromQLBuilder) BuildRate() string {
	base := b.buildBase()
	return fmt.Sprintf("rate(%s[%s])", base, formatDuration(b.timeWindow))
}

// BuildIRate builds an irate (instant rate) query with the specified time window.
func (b *PromQLBuilder) BuildIRate() string {
	base := b.buildBase()
	return fmt.Sprintf("irate(%s[%s])", base, formatDuration(b.timeWindow))
}

// BuildIncrease builds an increase query with the specified time window.
func (b *PromQLBuilder) BuildIncrease() string {
	base := b.buildBase()
	return fmt.Sprintf("increase(%s[%s])", base, formatDuration(b.timeWindow))
}

// BuildPercentile builds a histogram_quantile query for percentile calculation.
func (b *PromQLBuilder) BuildPercentile(percentile float64, bucketMetric string) string {
	base := b.buildBase()
	return fmt.Sprintf("histogram_quantile(%.2f, %s)", percentile/100, base)
}

// buildBase builds the base metric with labels (no aggregation).
func (b *PromQLBuilder) buildBase() string {
	query := b.metric
	if len(b.labels) > 0 {
		labels := make([]string, 0, len(b.labels))
		for k, v := range b.labels {
			labels = append(labels, fmt.Sprintf(`%s="%s"`, k, v))
		}
		query += "{" + strings.Join(labels, ",") + "}"
	}
	return query
}

// formatDuration formats a duration for PromQL (e.g., 5m, 1h, 1d).
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.0fh", d.Hours())
	}
	return fmt.Sprintf("%.0fd", d.Hours()/24)
}

// Helper functions for common query patterns.

// QueryRate builds a rate query for a metric.
func QueryRate(metric string, timeWindow time.Duration) string {
	return fmt.Sprintf("rate(%s[%s])", metric, formatDuration(timeWindow))
}

// QueryIRate builds an irate query for a metric.
func QueryIRate(metric string, timeWindow time.Duration) string {
	return fmt.Sprintf("irate(%s[%s])", metric, formatDuration(timeWindow))
}

// QueryAvgRate builds an avg(rate(...)) query.
func QueryAvgRate(metric string, labels map[string]string, timeWindow time.Duration) string {
	builder := NewPromQLBuilder(metric).
		WithLabels(labels).
		WithTimeWindow(timeWindow)
	return builder.BuildRate()
}

// QueryPercentile builds a quantile query for histogram metrics.
func QueryPercentile(metric string, percentile float64, labels map[string]string) string {
	builder := NewPromQLBuilder(metric).
		WithLabels(labels).
		WithPercentile(percentile)
	return builder.Build()
}

// QueryHTTPDurationRate queries HTTP request duration rate.
func QueryHTTPDurationRate(service, namespace string, timeWindow time.Duration) string {
	labels := make(map[string]string)
	if namespace != "" {
		labels["namespace"] = namespace
	}
	if service != "" {
		labels["service"] = service
	}
	return QueryAvgRate("http_request_duration_seconds_sum", labels, timeWindow)
}

// QueryHTTPRequestRate queries HTTP request rate.
func QueryHTTPRequestRate(service, namespace string, timeWindow time.Duration) string {
	labels := make(map[string]string)
	if namespace != "" {
		labels["namespace"] = namespace
	}
	if service != "" {
		labels["service"] = service
	}
	return QueryAvgRate("http_requests_total", labels, timeWindow)
}
