package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func sendTelegramMessage(message string, recipientID string) error {
	// Get Telegram API token from environment variables
	apiToken := os.Getenv("TELEGRAM_API_TOKEN")

	// Create URL for sending message
	sendMessageURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", apiToken)

	// Create JSON request body
	requestBody := map[string]interface{}{
		"chat_id": recipientID,
		"text":    message,
	}

	// Convert request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return err
	}

	// Create HTTP client and POST request
	client := http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest("POST", sendMessageURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	// Send request and handle response
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending HTTP request: %v", err)
		return err
	}
	defer response.Body.Close()

	// Check response status code
	if response.StatusCode != http.StatusOK {
		// Print response body message
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return err
		}
		var data struct {
			OK          bool   `json:"ok"`
			Description string `json:"description"`
		}
		err = json.Unmarshal(responseBody, &data)
		if err != nil {
			log.Printf("Error unmarshaling response: %v", err)
			return err
		}
		log.Printf("Telegram API error: unexpected status code %d, Description: %s", response.StatusCode, data.Description)

		return fmt.Errorf("Telegram API error: unexpected status code %d", response.StatusCode)
	}

	// Parse JSON response
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return err
	}
	var data struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return err
	}

	// Check if message was sent successfully
	if !data.OK {
		log.Printf("Telegram API error: %s", data.Description)
		return errors.New(data.Description)
	}

	return nil
}
