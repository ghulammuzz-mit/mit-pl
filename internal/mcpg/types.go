package mcpg

import (
	"time"
)

// Common types used across MCP tools.

// TimeRange represents a time range for queries.
type TimeRange struct {
	Start time.Time
	End   time.Time
	Step  time.Duration
}

// CommonTimeRanges returns predefined time ranges.
func CommonTimeRanges() map[string]TimeRange {
	now := time.Now()
	return map[string]TimeRange{
		"1h":  {Start: now.Add(-1 * time.Hour), End: now, Step: 1 * time.Minute},
		"6h":  {Start: now.Add(-6 * time.Hour), End: now, Step: 5 * time.Minute},
		"24h": {Start: now.Add(-24 * time.Hour), End: now, Step: 15 * time.Minute},
		"7d":  {Start: now.Add(-7 * 24 * time.Hour), End: now, Step: 1 * time.Hour},
		"30d": {Start: now.Add(-30 * 24 * time.Hour), End: now, Step: 6 * time.Hour},
	}
}

// ParseTimeRange parses a time range string.
// Supported formats: "1h", "6h", "24h", "7d", "30d"
func ParseTimeRange(timeRange string) (TimeRange, bool) {
	ranges := CommonTimeRanges()
	tr, ok := ranges[timeRange]
	if !ok {
		return TimeRange{}, false
	}
	return tr, true
}

// AggregationType represents the type of aggregation to apply.
type AggregationType string

const (
	AggregationAvg      AggregationType = "avg"
	AggregationSum      AggregationType = "sum"
	AggregationMax      AggregationType = "max"
	AggregationMin      AggregationType = "min"
	AggregationCount    AggregationType = "count"
	AggregationPercentile AggregationType = "percentile"
)

// QueryResult represents a simplified query result.
type QueryResult struct {
	Metric map[string]string
	Values []TimeValue
}

// TimeValue represents a timestamp-value pair.
type TimeValue struct {
	Timestamp int64
	Value     float64
}

// MetricFilter represents filters for metrics.
type MetricFilter struct {
	Namespace string
	Service   string
	Pod       string
	Labels    map[string]string
}
