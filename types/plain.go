package types

import (
	"io/ioutil"
	"net/http"
)

func ParsePlainBody(r *http.Request) (string, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "Failed to read request body", err
	}

	if len(bodyBytes) == 0 {
		return "Empty plain text", nil
	}

	return string(bodyBytes), nil
}
