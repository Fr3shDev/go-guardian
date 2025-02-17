package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/common-nighthawk/go-figure"
)

type WebsiteStatus struct {
	URL              string
	ResponseTime     time.Duration
	StatusCode       int
	SSLExpiration    time.Time
	SSLExpiryWarning bool
}

func checkWebsite(url string, sslExpiryThreshold time.Duration) (*WebsiteStatus, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()
	response, error := client.Get(url)
	if error != nil {
		return nil, error
	}
	defer response.Body.Close()

	duration := time.Since(start)

	status := &WebsiteStatus{
		URL:          url,
		ResponseTime: duration,
		StatusCode:   response.StatusCode,
	}

	// if HTTPS, check the SSL certificate details
	if response.Request.URL.Scheme == "https" {
		// Establish a TLS connection to retrieve certificate info.
		connection, error := tls.Dial("tcp", response.Request.URL.Host+":443", &tls.Config{
			// Verify the certificate properly in production
			// Here we use InsecureSkipVerify to simply fetch the certificate details
			InsecureSkipVerify: true,
		})
		if error != nil {
			return status, error
		}
		defer connection.Close()
		certificates := connection.ConnectionState().PeerCertificates
		if len(certificates) > 0 {
			certificate := certificates[0]
			status.SSLExpiration = certificate.NotAfter
			if time.Until(certificate.NotAfter) < sslExpiryThreshold {
				status.SSLExpiryWarning = true
			}
		}
	}
	return status, nil
}

func checkWebsitesConcurrently(websites []string, sslExpiryThreshold time.Duration) {
	var wg sync.WaitGroup
	for _, url := range websites {
		wg.Add(1)
		// Launch a goroutine for each website.
		go func(site string) {
			defer wg.Done()
			status, error := checkWebsite(site, sslExpiryThreshold)
			if error != nil {
				fmt.Printf("Error checking %s: %v\n", site, error)
				return
			}
			fmt.Printf("URL: %s | Status Code: %d | Response Time: %v\n", status.URL, status.StatusCode, status.ResponseTime)
			if !status.SSLExpiration.IsZero() {
				fmt.Printf(" SSL Certificate expires on: %s\n", status.SSLExpiration.Format("2006-01-02"))
				if status.SSLExpiryWarning {
					fmt.Printf(" WARNING: SSL certificate expires soon!\n")
				}
			}
		}(url)
	}
	wg.Wait()
}

func main() {
	// List of websites to monitor.
	websites := []string{
		// "https://www.google.com",
		// "https://www.github.com",
		// "http://example.com",
		"https://techchantier.com",
	}

	// Define a threshold for SSL certificate expiration (e.g., 30days).
	sslExpiryThreshold := 30 * 24 * time.Hour

	// Create a ticker to schedule checks every minute.
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	figure.NewColorFigure("Starting website health monitoring...", "", "green", true).Print()
	// myFigure.Print()
	// fmt.Println("Starting website health monitoring...")
	// Perform an initial check.
	checkWebsitesConcurrently(websites, sslExpiryThreshold)

	// Run periodic checks.
	for range ticker.C {
		fmt.Println("\nPerforming scheduled check...")
		checkWebsitesConcurrently(websites, sslExpiryThreshold)
	}
}
