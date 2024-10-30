package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Prompt struct {
	Prompt string
	Width  int16
	Height int16
}

type BlackForest struct {
	Url    string
	ApiKey string
	Prompt map[string]interface{}
}

func (bf *BlackForest) Init() {
	// Define the API URL
	bf.Url = "https://api.bfl.ml/v1/flux-pro-1.1"

	// API key from environment variable
	bf.ApiKey = os.Getenv("BFL_API_KEY")
	//bf.ApiKey = "83fd74e4-5612-4430-9b16-cc430a4473d8"

}

func (bf *BlackForest) SetPrompt(prompt string) {
	// Define the request payload
	bf.Prompt = map[string]interface{}{
		// "prompt": "A cat on its back legs running like a human is holding a big silver fish with its arms. The cat is running away from the shop owner and has a panicked look on his face. The scene is situated in a crowded market.",
		"prompt": "You are a creative assistant, who seeks to guide and help a marketer to create content on social media, create an image for the next prompt" + prompt,
		"width":  1024,
		"height": 768,
	}

}

func (bf *BlackForest) Request() string {

	// Convert the payload to JSON
	jsonData, err := json.Marshal(bf.Prompt)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}

	// Create a new request
	req, err := http.NewRequest("POST", bf.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Add headers
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-key", bf.ApiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Print the response
	fmt.Println("Response:", string(body))

	// Extract request_id using jq-like functionality in Go
	var result map[string]interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Request ID: %v\n", result["id"])
	requestID, ok := result["id"].(string)
	fmt.Printf("Request ID: %v\n", requestID)
	if !ok {
		log.Fatalf("Error: 'id' field is not a string or is missing")
	}

	return requestID
}

func (bf *BlackForest) Poll(requestID string) string {

	// Request ID obtained from previous step (replace this with the actual request_id or set dynamically)
	// requestID := "your-request-id-here"

	// Base URL for checking result
	url := fmt.Sprintf("https://api.bfl.ml/v1/get_result?id=%s", requestID)

	// Infinite loop to check the status
	for {
		// Sleep for 0.5 seconds between requests
		time.Sleep(500 * time.Millisecond)

		// Create the request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		// Add headers
		req.Header.Set("accept", "application/json")
		req.Header.Set("x-key", bf.ApiKey)

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}

		// Parse the JSON response
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatalf("Error unmarshalling JSON: %v", err)
		}

		// Extract the status
		status, ok := result["status"].(string)
		if !ok {
			log.Fatalf("Error retrieving status from response")
		}

		// Check if the status is "Ready"
		if status == "Ready" {
			// Extract the result
			resultSample, ok := result["result"].(map[string]interface{})
			if !ok {
				log.Fatalf("Error retrieving result from response")
			}
			sample, ok := resultSample["sample"].(string)
			if !ok {
				log.Fatalf("Error retrieving sample from result")
			}
			// Print the result and break the loop
			fmt.Printf("Result: %s\n", sample)
			return sample
			break
		} else {
			// Print the current status
			fmt.Printf("Status: %s\n", status)
		}
	}
	return ""
}
