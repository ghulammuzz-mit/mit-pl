package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	grafanaURL := os.Getenv("GRAFANA_URL")
	grafanaAPIKey := os.Getenv("GRAFANA_API_KEY")
	if grafanaURL == "" || grafanaAPIKey == "" {
		log.Fatal("GRAFANA_URL and GRAFANA_API_KEY must be set (via .env or environment)")
	}

	cmd := exec.Command("go", "run", "./cmd/mcp-grafana")
	cmd.Env = append(os.Environ(),
		"GRAFANA_URL="+grafanaURL,
		"GRAFANA_API_KEY="+grafanaAPIKey,
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Printf("[LOG] %s\n", scanner.Text())
		}
	}()

	go func() {
		defer wg.Done()
		decoder := json.NewDecoder(stdout)
		for {
			var response map[string]interface{}
			if err := decoder.Decode(&response); err != nil {
				break
			}
			jsonBytes, _ := json.MarshalIndent(response, "", "  ")
			fmt.Printf("[RESPONSE]\n%s\n", string(jsonBytes))
		}
	}()

	initReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	encoder := json.NewEncoder(stdin)
	if err := encoder.Encode(initReq); err != nil {
		log.Fatal(err)
	}

	notif := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "notifications/initialized",
	}
	if err := encoder.Encode(notif); err != nil {
		log.Fatal(err)
	}

	listReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/list",
	}
	if err := encoder.Encode(listReq); err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	stdin.Close()
	wg.Wait()
	cmd.Wait()
}
