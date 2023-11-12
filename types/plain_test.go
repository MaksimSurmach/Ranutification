package types

import (
	"net/http"
	"strings"
	"testing"
)

func TestParsePlainBody(t *testing.T) {
	body := "Plain text content"
	req, err := http.NewRequest("POST", "/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	result, err := ParsePlainBody(req)
	if err != nil {
		t.Errorf("Error in ParsePlainBody: %v", err)
	}
	if result != "Plain text content" {
		t.Errorf("Expected: Plain text content, Got: %s", result)
	}
}

func TestParsePlainBodyEmpty(t *testing.T) {
	body := ""
	req, err := http.NewRequest("POST", "/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	result, err := ParsePlainBody(req)
	if err != nil {
		t.Errorf("Error in ParsePlainBody: %v", err)
	}
	if result != "Empty plain text" {
		t.Errorf("Expected: Empty plain text, Got: %s", result)
	}
}

func TestParsePlainBodyInvalid(t *testing.T) {
	body := "This is plain text"
	req, err := http.NewRequest("POST", "/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	_, err = ParsePlainBody(req)
	if err != nil {
		t.Errorf("Expected no error in ParsePlainBody, Got: %v", err)
	}
}
