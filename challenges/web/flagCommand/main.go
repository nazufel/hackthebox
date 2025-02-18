package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// getRequest makes a GET request to the specified URL and returns the response body as a string.
func makeRequest(url string, payload map[string]string) (string, error) {
	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Create io.Reader from JSON bytes
	payloadReader := bytes.NewReader(jsonData)

	// Make the POST request with the JSON reader
	resp, err := http.Post(url, "application/json", payloadReader)
	if err != nil {
		return "", fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

func main() {

	// Get URL from command line args
	if len(os.Args) < 2 {
		fmt.Println("Error: URL argument is required")
		os.Exit(1)
	}
	url := os.Args[1] + "/api/monitor"

	payload := map[string]string{
		"command": "Blip-blop, in a pickle with a hiccup! Shmiggity-shmack",
	}

	// Make the request
	response, err := makeRequest(url, payload)
	if err != nil {
		fmt.Printf("Error making GET request: %v\n", err)
		return
	}

	// Print the response
	fmt.Println(response)
}
