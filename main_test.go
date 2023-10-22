package main

import (
	"io/ioutil"
	"net/http"
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
			t.Logf("Request: %v", req)
		}
		if result != test.expected {
			t.Errorf("Expected: %s, Got: %s", test.expected, result)
		}
	}
}
