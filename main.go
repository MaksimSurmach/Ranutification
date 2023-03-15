package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func sendTelegramMessage(message string) error {
	// Get Telegram API token and recipient ID from environment variables
	apiToken := os.Getenv("TELEGRAM_API_TOKEN")
	recipientID := os.Getenv("TELEGRAM_RECIPIENT_ID")

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
		return err
	}

	// Create HTTP client and POST request
	client := http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest("POST", sendMessageURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	// Send request and handle response
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check response status code
	if response.StatusCode != http.StatusOK {
		// print response body mesasge
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		var data struct {
			OK          bool   `json:"ok"`
			Description string `json:"description"`
		}
		err = json.Unmarshal(responseBody, &data)
		if err != nil {
			return err
		}
		print(data.Description)

		return fmt.Errorf("Telegram API error: unexpected status code %d", response.StatusCode)
	}

	// Parse JSON response
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var data struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return err
	}

	// Check if message was sent successfully
	if !data.OK {
		return errors.New(data.Description)
	}

	return nil
}

func parseBodyJson(r *http.Request) (string, error) {
	decoder := json.NewDecoder(r.Body)
	var t map[string]interface{}
	err := decoder.Decode(&t)
	if err != nil {
		return "", err
	}
	var body string
	for k, v := range t {
		body += fmt.Sprintf("%s: %v\n", k, v)
	}
	return body, nil
}

func parseBody(r *http.Request) (string, error) {
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		body, err := parseBodyJson(r)
		if err != nil {
			return "Invalid JSON body", err
		}
		fmt.Println(body)
		return body, nil
	case "text/plain":
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return "Failed to read request body", err
		}
		body := string(bodyBytes)
		fmt.Println(body)
		return body, nil
	default:
		return "Invalid content type", errors.New("Invalid content type")
	}
}

// handler for /health
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

// handler for webhook
func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	CODE := os.Getenv("CODE")
	//CODE := "test"
	if r.URL.Path != "/"+CODE {
		http.Error(w, "Invalid code", http.StatusBadRequest)
		return
	}

	// Parse incoming data to string
	body, err := parseBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if body is empty after parsing
	if body == "" {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}

	// Get who sent the message
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	body += "\nSender: " + IPAddress + "\n"

	err = sendTelegramMessage(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response 200 OK
	fmt.Fprintf(w, "OK")

}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(":8080", nil)

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		fmt.Printf("Received signal: %v\n", sig)
		sendTelegramMessage("Bot stopped")
		os.Exit(1)
	}()

	// Handle internal panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
			sendTelegramMessage("Bot stopped")
			os.Exit(1)
		}
	}()

}
