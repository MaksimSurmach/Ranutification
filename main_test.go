package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseBody(t *testing.T) {
	tests := []struct {
		contentType string
		body        string
		expected    string
	}{
		{
			contentType: "application/json",
			body:        `{"key": "value"}`,
			expected:    "key: value\n",
		},
		{
			contentType: "text/plain",
			body:        "Plain text content",
			expected:    "Plain text content",
		},
		{
			contentType: "application/xml",
			body:        "<root><element>Value</element></root>",
			expected:    "Value",
		},
		{
			contentType: "invalid-content-type",
			body:        "Invalid content",
			expected:    "Invalid content type",
		},
		{
			contentType: "",
			body:        `{"key": "value"}`,
			expected:    "key: value\n",
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest("POST", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		if test.contentType != "" {
			req.Header.Set("Content-Type", test.contentType)
		}
		req.Body = ioutil.NopCloser(strings.NewReader(test.body))
		result, err := parseBody(req)
		if err != nil {
			t.Errorf("Error in parseBody: %v", err)
		}
		if result != test.expected {
			t.Errorf("Expected: %s, Got: %s", test.expected, result)
		}
	}
}

func TestMainHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/123456789", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	mainHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, Got %d", http.StatusOK, w.Code)
	}
}
