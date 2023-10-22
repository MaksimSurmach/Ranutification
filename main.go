package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"ranutification/types"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func parseBody(r *http.Request) (string, error) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "" {
		body, err := types.ParseBodyJson(r)
		if err == nil {
			return body, nil
		}

		body, err = types.ParseXMLBody(r)
		if err == nil {
			return body, nil
		}

		body, err = types.ParsePlainBody(r)
		if err != nil {
			return "Invalid content type", errors.New("Invalid content type")
		}
		return body, nil
	}

	switch contentType {
	case "application/json":
		return types.ParseBodyJson(r)
	case "text/plain":
		return types.ParsePlainBody(r)
	case "application/xml":
		return types.ParseXMLBody(r)
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
	// Parse incoming data to string
	body, err := parseBody(r)
	if err != nil {
		log.Printf("Error parsing body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if body is empty after parsing
	if body == "" {
		log.Println("Empty body")
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

	// Используйте имя функции sendTelegramMessage из telegram.go
	err = sendTelegramMessage(body, r.URL.Path[1:])
	if err != nil {
		log.Printf("Error sending Telegram message: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response 200 OK
	w.Write([]byte("OK"))
}

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", healthHandler)

	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(3, 1*time.Minute))
		r.Post("/{^[0-9]{6,9}$}", mainHandler)
	})

	log.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", r)
}
