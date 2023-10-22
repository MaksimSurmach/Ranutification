package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseJSONBody(r *http.Request) (string, error) {
	decoder := json.NewDecoder(r.Body)
	var t map[string]interface{}
	err := decoder.Decode(&t)
	if err != nil {
		return "", err
	}
	var body string

	if len(t) == 0 {
		return "Empty JSON object", nil
	}

	for k, v := range t {
		body += k + ": " + fmt.Sprint(v) + "\n"
	}
	return body, nil
}
