package types

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type XMLBody struct {
	XMLName xml.Name
	Content string `xml:",chardata"`
}

func ParseXMLBody(r *http.Request) (string, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "Failed to read request body", err
	}

	var xmlBody XMLBody
	err = xml.Unmarshal(bodyBytes, &xmlBody)
	if err != nil {
		return "Failed to parse XML body", err
	}

	return xmlBody.Content, nil
}
