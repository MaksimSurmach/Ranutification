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

	return string(bodyBytes), nil
}
