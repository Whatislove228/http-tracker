package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Result struct {
	URL        string
	Status     string
	StatusCode int
	Duration   time.Duration
	Error      error
}

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorBold  = "\033[1m"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: uptime-checker <sites-file>")
		return
	}

	file := os.Args[1]

	sites, err := readSites(file)

	if err != nil {
		fmt.Println("Error reading file: ", file)
		return
	}

	results := make(chan Result)
	for _, url := range sites {
		go func(siteUrl string) {
			results <- checkSite(siteUrl)
		}(url)
	}

	fmt.Printf(colorBold+"%-30s %-5s %-8s %s\n"+colorReset, "URL", "STATE", "CODE", "DETAIL")
	fmt.Println(strings.Repeat("-", 70))

	for range sites {
		result := <-results

		printResult(result)

	}
}

func coloredStatus(status string) string {
	if status == "UP" {
		return colorGreen + status + colorReset
	}

	if status == "DOWN" {
		return colorRed + status + colorReset
	}

	return status
}

func printResult(result Result) {
	if result.Error != nil {
		fmt.Printf("%-30s %-14s %-8s %v\n", result.URL, coloredStatus(result.Status), "-", result.Error)
		return
	}

	fmt.Printf("%-30s %-14s %-8d %v\n", result.URL, coloredStatus(result.Status), result.StatusCode, result.Duration)
}

func readSites(file string) ([]string, error) {
	data, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	var sites []string
	for _, line := range lines {
		site := strings.TrimSpace(line)

		if site == "" {
			continue
		}

		if strings.HasPrefix(site, "#") {
			continue
		}

		sites = append(sites, site)
	}

	return sites, nil
}

func checkSite(url string) Result {
	start := time.Now()

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	duration := time.Since(start)

	if err != nil {
		return Result{
			URL:      url,
			Status:   "DOWN",
			Duration: duration,
			Error:    err,
		}
	}

	status := "DOWN"
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		status = "UP"
	}

	result := Result{
		URL:        url,
		Status:     status,
		StatusCode: resp.StatusCode,
		Duration:   duration,
	}

	resp.Body.Close()

	return result
}
