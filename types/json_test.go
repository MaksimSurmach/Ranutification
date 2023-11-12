package types

import (
	"net/http"
	"strings"
	"testing"
)

func TestParseJSONBody(t *testing.T) {
	body := `{"key": "value"}`
	req, err := http.NewRequest("POST", "/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	result, err := ParseJSONBody(req)
	if err != nil {
		t.Errorf("Error in ParseJSONBody: %v", err)
	}
	if result != "key: value\n" {
		t.Errorf("Expected: key: value, Got: %s", result)
	}
}

func TestParseJSONBodyEmpty(t *testing.T) {
	body := `{}`
	req, err := http.NewRequest("POST", "/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	result, err := ParseJSONBody(req)
	if err != nil {
		t.Errorf("Error in ParseJSONBody: %v", err)
	}
	if result != "Empty JSON object" {
		t.Errorf("Expected: Empty JSON object, Got: %s", result)
	}
}

func TestParseJSONBodyInvalid(t *testing.T) {
	body := `invalid-json`
	req, err := http.NewRequest("POST", "/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	_, err = ParseJSONBody(req)
	if err == nil {
		t.Error("Expected error in ParseJSONBody for invalid JSON")
	}
}
