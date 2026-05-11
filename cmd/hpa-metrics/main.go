package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"mit/platform/internal/grafana"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type ServiceMetrics struct {
	Name            string
	CPUAvg          float64
	CPUMax          float64
	CPUP95          float64
	MemoryAvg        float64
	MemoryMax        float64
	MemoryP95       float64
	CurrentReplicas  int
}

type HPARecommendation struct {
	Service           string
	RecommendedCPU    CPURecommendation
	RecommendedMemory MemoryRecommendation
	MinReplicas      int
	MaxReplicas      int
	TargetCPUUtil    int
	TargetMemUtil    int
	Notes            string
}

type CPURecommendation struct {
	Millis     string
	Cores      string
	RawCores   float64
}

type MemoryRecommendation struct {
	MiB    string
	GiB     string
	Bytes   int64
}

func main() {
	// Load environment
	if err := loadEnv(); err != nil {
		slog.Error("Failed to load env", "error", err)
		os.Exit(1)
	}

	client := grafana.New()

	services := []string{
		"lvservices", "s-backend", "s-dummy-p1", "s-dummy-pw1",
		"s-hook", "s-multibackend", "s-web", "s-backendapi", "s-openapi",
	}

	fmt.Println("🔍 Collecting metrics for HPA planning...")
	fmt.Println()

	allMetrics := make([]ServiceMetrics, 0, len(services))
	ctx := context.Background()
	timeRange := 5 * 24 * time.Hour
	end := time.Now()
	start := end.Add(-timeRange)

	for _, svc := range services {
		metrics, err := collectServiceMetrics(ctx, client, svc, start, end)
		if err != nil {
			slog.Warn("Failed to collect metrics for service", "service", svc, "error", err)
			continue
		}
		allMetrics = append(allMetrics, *metrics)
		fmt.Printf("✓ %s: CPU avg=%.4f cores, Memory avg=%.1f MiB\n",
			svc, metrics.CPUAvg, metrics.MemoryAvg/(1024*1024))
	}

	fmt.Println()
	fmt.Println("📊 Generating HPA recommendations...")
	fmt.Println()

	// Create hpa-planning directory
	if err := os.MkdirAll("hpa-planning", 0755); err != nil {
		slog.Error("Failed to create directory", "error", err)
		os.Exit(1)
	}

	// Generate recommendations
	recommendations := generateRecommendations(allMetrics)

	// Write markdown report
	if err := writeMarkdownReport(recommendations, allMetrics); err != nil {
		slog.Error("Failed to write report", "error", err)
		os.Exit(1)
	}

	// Write JSON output
	if err := writeJSONReport(allMetrics, recommendations); err != nil {
		slog.Error("Failed to write JSON", "error", err)
		os.Exit(1)
	}

	fmt.Println("✅ Report generated in hpa-planning/ folder")
	fmt.Println("   - hpa-planning/report.md")
	fmt.Println("   - hpa-planning/metrics.json")
}

func collectServiceMetrics(ctx context.Context, client *grafana.Client, service string, start, end time.Time) (*ServiceMetrics, error) {
	// Query CPU (rate over 5m windows, aggregated)
	cpuQuery := fmt.Sprintf(`avg(rate(container_cpu_usage_seconds_total{namespace="prod-app", pod=~"%s-.*"}[5m]))`, service)
	cpuResult, err := client.DirectQuery(ctx, "prometheus", cpuQuery, start, end, 1*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("CPU query failed: %w", err)
	}

	// Query Memory
	memQuery := fmt.Sprintf(`avg(container_memory_working_set_bytes{namespace="prod-app", pod=~"%s-.*"})`, service)
	memResult, err := client.DirectQuery(ctx, "prometheus", memQuery, start, end, 1*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("memory query failed: %w", err)
	}

	metrics := &ServiceMetrics{
		Name: service,
	}

	// Process CPU data
	if len(cpuResult.Data.Result) > 0 && len(cpuResult.Data.Result[0].Values) > 0 {
		values := extractFloatValues(cpuResult.Data.Result[0].Values)
		metrics.CPUAvg = avg(values)
		metrics.CPUMax = maxVal(values)
		metrics.CPUP95 = percentileVal(values, 95)
	}

	// Process Memory data
	if len(memResult.Data.Result) > 0 && len(memResult.Data.Result[0].Values) > 0 {
		values := extractFloatValues(memResult.Data.Result[0].Values)
		metrics.MemoryAvg = avg(values)
		metrics.MemoryMax = maxVal(values)
		metrics.MemoryP95 = percentileVal(values, 95)
	}

	return metrics, nil
}

func extractFloatValues(values [][]any) []float64 {
	result := make([]float64, 0, len(values))
	for _, v := range values {
		if len(v) >= 2 {
			if valStr, ok := v[1].(string); ok {
				var f float64
				fmt.Sscanf(valStr, "%f", &f)
				result = append(result, f)
			}
		}
	}
	return result
}

func avg(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func maxVal(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for _, v := range values {
		if v > m {
			m = v
		}
	}
	return m
}

func percentileVal(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)

	// Simple sort
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	idx := (p / 100) * float64(len(sorted)-1)
	lower := int(idx)
	upper := lower + 1
	if upper >= len(sorted) {
		return sorted[len(sorted)-1]
	}
	weight := idx - float64(lower)
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}

func generateRecommendations(metrics []ServiceMetrics) []HPARecommendation {
	recs := make([]HPARecommendation, 0, len(metrics))

	for _, m := range metrics {
		rec := HPARecommendation{
			Service: m.Name,
		}

		// CPU Recommendation: Request = avg * 1.5, Limit = max * 1.2
		cpuRequestCores := m.CPUAvg * 1.5
		cpuLimitCores := math.Max(m.CPUMax*1.2, cpuRequestCores*1.5)

		// Ensure minimum values
		cpuRequestCores = math.Max(cpuRequestCores, 0.01)  // 10m minimum
		cpuLimitCores = math.Max(cpuLimitCores, 0.05)     // 50m minimum

		rec.RecommendedCPU = CPURecommendation{
			Cores:    fmt.Sprintf("%.2f", cpuLimitCores),
			Millis:    fmt.Sprintf("%.0fm", cpuLimitCores*1000),
			RawCores: cpuLimitCores,
		}

		// Memory Recommendation: Request = avg * 1.5, Limit = p95 * 1.2
		memRequestBytes := int64(math.Max(m.MemoryAvg*1.5, 128*1024*1024))  // 128Mi minimum
		memLimitBytes := int64(math.Max(m.MemoryP95*1.2, float64(memRequestBytes)*1.5))

		rec.RecommendedMemory = MemoryRecommendation{
			Bytes: memLimitBytes,
			MiB:   fmt.Sprintf("%.0f", float64(memLimitBytes)/(1024*1024)),
			GiB:    fmt.Sprintf("%.2f", float64(memLimitBytes)/(1024*1024*1024)),
		}

		// HPA Settings based on CPU utilization
		if m.CPUAvg < 0.01 {
			rec.MinReplicas = 1
			rec.MaxReplicas = 2
			rec.TargetCPUUtil = 80
			rec.TargetMemUtil = 80
			rec.Notes = "Low utilization - minimal scaling"
		} else if m.CPUAvg < 0.1 {
			rec.MinReplicas = 2
			rec.MaxReplicas = 5
			rec.TargetCPUUtil = 70
			rec.TargetMemUtil = 75
			rec.Notes = "Low utilization - conservative scaling"
		} else if m.CPUAvg < 0.5 {
			rec.MinReplicas = 2
			rec.MaxReplicas = 10
			rec.TargetCPUUtil = 70
			rec.TargetMemUtil = 75
			rec.Notes = "Medium utilization - normal scaling"
		} else {
			rec.MinReplicas = 3
			rec.MaxReplicas = 20
			rec.TargetCPUUtil = 60
			rec.TargetMemUtil = 70
			rec.Notes = "High utilization - aggressive scaling"
		}

		recs = append(recs, rec)
	}

	return recs
}

func writeMarkdownReport(recs []HPARecommendation, metrics []ServiceMetrics) error {
	f, err := os.Create("hpa-planning/report.md")
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "# HPA Planning Report")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "**Generated:**", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintln(f, "**Data Period:** Last 5 days")
	fmt.Fprintln(f, "**Namespace:** prod-app")
	fmt.Fprintln(f, "")

	fmt.Fprintln(f, "## Summary")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "| Service | CPU Avg | CPU Max | Memory Avg | Memory Max | Min | Max | Target CPU | Target Mem | Notes |")
	fmt.Fprintln(f, "|----------|----------|----------|------------|------------|-----|-----|------------|-------------|--------|")

	for i, rec := range recs {
		m := metrics[i]
		fmt.Fprintf(f, "| %s | %.4f | %.4f | %.1f MiB | %.1f MiB | %d | %d | %d%% | %d%% | %s |\n",
			rec.Service,
			m.CPUAvg,
			m.CPUMax,
			m.MemoryAvg/(1024*1024),
			m.MemoryMax/(1024*1024),
			rec.MinReplicas,
			rec.MaxReplicas,
			rec.TargetCPUUtil,
			rec.TargetMemUtil,
			rec.Notes)
	}

	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "## Detailed Recommendations")
	fmt.Fprintln(f, "")

	for _, rec := range recs {
		fmt.Fprintln(f, "###", rec.Service)
		fmt.Fprintln(f, "")
		fmt.Fprintln(f, "**Resource Requests:**")
		fmt.Fprintf(f, "- CPU: `%sm` (request)\n", rec.RecommendedCPU.Millis)
		fmt.Fprintf(f, "- Memory: `%sMi` (request)\n", rec.RecommendedMemory.MiB)
		fmt.Fprintln(f, "")
		fmt.Fprintln(f, "**Resource Limits:**")
		fmt.Fprintf(f, "- CPU: `%s` (limit)\n", rec.RecommendedCPU.Cores)
		fmt.Fprintf(f, "- Memory: `%sMi` (limit)\n", rec.RecommendedMemory.MiB)
		fmt.Fprintln(f, "")
		fmt.Fprintln(f, "**HPA Configuration:**")
		fmt.Fprintf(f, "- Min Replicas: `%d`\n", rec.MinReplicas)
		fmt.Fprintf(f, "- Max Replicas: `%d`\n", rec.MaxReplicas)
		fmt.Fprintf(f, "- Target CPU Utilization: `%d%%`\n", rec.TargetCPUUtil)
		fmt.Fprintf(f, "- Target Memory Utilization: `%d%%`\n", rec.TargetMemUtil)
		fmt.Fprintln(f, "")
		fmt.Fprintln(f, "**Notes:**", rec.Notes)
		fmt.Fprintln(f, "")

		// HPA YAML example
		fmt.Fprintln(f, "```yaml")
		fmt.Fprintln(f, "apiVersion: autoscaling/v2")
		fmt.Fprintln(f, "kind: HorizontalPodAutoscaler")
		fmt.Fprintln(f, "metadata:")
		fmt.Fprintf(f, "  name: %s-hpa\n", rec.Service)
		fmt.Fprintln(f, "spec:")
		fmt.Fprintln(f, "  scaleTargetRef:")
		fmt.Fprintln(f, "    apiVersion: apps/v1")
		fmt.Fprintln(f, "    kind: Deployment")
		fmt.Fprintf(f, "    name: %s\n", rec.Service)
		fmt.Fprintln(f, "  minReplicas:", rec.MinReplicas)
		fmt.Fprintln(f, "  maxReplicas:", rec.MaxReplicas)
		fmt.Fprintln(f, "  metrics:")
		fmt.Fprintf(f, "  - type: Resource\n")
		fmt.Fprintf(f, "    resource:\n")
		fmt.Fprintf(f, "      name: cpu\n")
		fmt.Fprintf(f, "      target:\n")
		fmt.Fprintf(f, "        type: Utilization\n")
		fmt.Fprintf(f, "        averageUtilization: %d\n", rec.TargetCPUUtil)
		fmt.Fprintf(f, "  - type: Resource\n")
		fmt.Fprintf(f, "    resource:\n")
		fmt.Fprintf(f, "      name: memory\n")
		fmt.Fprintf(f, "      target:\n")
		fmt.Fprintf(f, "        type: Utilization\n")
		fmt.Fprintf(f, "        averageUtilization: %d\n", rec.TargetMemUtil)
		fmt.Fprintln(f, "```")
		fmt.Fprintln(f, "")
	}

	return nil
}

func writeJSONReport(metrics []ServiceMetrics, recs []HPARecommendation) error {
	data := map[string]any{
		"generated_at": time.Now().Format(time.RFC3339),
		"period":      "5d",
		"namespace":   "prod-app",
		"metrics":     metrics,
		"recommendations": recs,
	}

	f, err := os.Create("hpa-planning/metrics.json")
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func loadEnv() error {
	// Try to load .env file
	for _, path := range []string{".env", "../.env", "../../.env"} {
		if _, err := os.Stat(path); err == nil {
			return godotenv.Load(path)
		}
	}
	return nil // No .env file found, that's OK
}
