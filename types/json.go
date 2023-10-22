package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseBodyJson(r *http.Request) (string, error) {
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
